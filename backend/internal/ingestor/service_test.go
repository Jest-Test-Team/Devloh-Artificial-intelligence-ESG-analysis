package ingestor

import (
	"encoding/json"
	"testing"

	"devlog-ai/internal/api/models"
	"devlog-ai/internal/privacy"
)

// 測試用的模擬Kafka生產者
type testKafkaProducer struct {
	messages []models.LogEntry
}

func (p *testKafkaProducer) ProduceMessage(topic string, key string, value []byte) error {
	var logEntry models.LogEntry
	json.Unmarshal(value, &logEntry)
	p.messages = append(p.messages, logEntry)
	return nil
}

func (p *testKafkaProducer) SendLogBatch(logs []models.LogEntry) error {
	p.messages = append(p.messages, logs...)
	return nil
}

func (p *testKafkaProducer) Close() error {
	return nil
}

func TestIngestorService_IngestLog(t *testing.T) {
	// 創建模擬Kafka生產者
	producer := &testKafkaProducer{
		messages: make([]models.LogEntry, 0),
	}

	// 創建PII檢測器
	piiDetector := privacy.NewPIIDetector()

	// 創建Ingestor服務
	service := NewService(producer, piiDetector)

	// 測試日誌處理
	logEntry := models.LogEntry{
		Service:   "test-service",
		Level:     "INFO",
		Message:   "User email: user@example.com and phone: 123-456-7890",
		Timestamp: "2023-07-15T12:00:00Z",
		UserID:    "test-user-123",
		Metadata: map[string]interface{}{
			"request_id": "req-123",
		},
	}

	// 處理日誌
	err := service.IngestLog(logEntry)
	if err != nil {
		t.Fatalf("Failed to ingest log: %v", err)
	}

	// 檢查是否有一條消息被發送
	if len(producer.messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(producer.messages))
	}

	// 檢查PII是否被匿名化
	processedMsg := producer.messages[0]
	if processedMsg.Message == logEntry.Message {
		t.Error("PII was not anonymized in the message")
	}

	// 檢查是否包含[EMAIL]和[PHONE]標記
	if !contains(processedMsg.Message, "[EMAIL]") {
		t.Error("Email was not anonymized properly")
	}
	if !contains(processedMsg.Message, "[PHONE]") {
		t.Error("Phone number was not anonymized properly")
	}

	// 檢查其他字段是否保持不變
	if processedMsg.Service != logEntry.Service {
		t.Errorf("Service field was modified: expected %s, got %s", logEntry.Service, processedMsg.Service)
	}
	if processedMsg.Level != logEntry.Level {
		t.Errorf("Level field was modified: expected %s, got %s", logEntry.Level, processedMsg.Level)
	}
	if processedMsg.UserID != logEntry.UserID {
		t.Errorf("UserID field was modified: expected %s, got %s", logEntry.UserID, processedMsg.UserID)
	}
}

func TestIngestorService_IngestLogBatch(t *testing.T) {
	// 創建模擬Kafka生產者
	producer := &testKafkaProducer{
		messages: make([]models.LogEntry, 0),
	}

	// 創建PII檢測器
	piiDetector := privacy.NewPIIDetector()

	// 創建Ingestor服務
	service := NewService(producer, piiDetector)

	// 測試批量日誌處理
	logEntries := []models.LogEntry{
		{
			Service:   "test-service-1",
			Level:     "INFO",
			Message:   "Credit card: 4111-1111-1111-1111",
			Timestamp: "2023-07-15T12:00:00Z",
			UserID:    "test-user-123",
		},
		{
			Service:   "test-service-2",
			Level:     "WARN",
			Message:   "IP address: 192.168.1.1",
			Timestamp: "2023-07-15T12:01:00Z",
			UserID:    "test-user-123",
		},
	}

	// 處理日誌批次
	err := service.IngestLogBatch(logEntries)
	if err != nil {
		t.Fatalf("Failed to ingest log batch: %v", err)
	}

	// 檢查是否有兩條消息被發送
	if len(producer.messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(producer.messages))
	}

	// 檢查每條消息的PII是否被匿名化
	if !contains(producer.messages[0].Message, "[CREDIT_CARD]") {
		t.Error("Credit card was not anonymized properly")
	}
	if !contains(producer.messages[1].Message, "[IP_ADDRESS]") {
		t.Error("IP address was not anonymized properly")
	}
}

// 輔助函數：檢查字符串是否包含子字符串
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
