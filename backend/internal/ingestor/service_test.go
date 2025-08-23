package ingestor

import (
	"context"
	"encoding/json"
	"testing"

	"devlog-ai/internal/api/models"
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
	// 創建Ingestor服務
	service := NewIngestorService()

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

	// 創建上下文
	ctx := context.Background()

	// 處理日誌
	err := service.IngestLog(ctx, logEntry)
	if err != nil {
		t.Fatalf("Failed to ingest log: %v", err)
	}

	// 由於我們沒有真正的Kafka生產者，這裡只測試函數是否執行成功
	// 在實際測試中，可以使用模擬的Kafka生產者來驗證消息發送
}

func TestIngestorService_IngestLogBatch(t *testing.T) {
	// 創建Ingestor服務
	service := NewIngestorService()

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

	// 創建上下文
	ctx := context.Background()

	// 處理日誌批次
	err := service.IngestLogBatch(ctx, logEntries)
	if err != nil {
		t.Fatalf("Failed to ingest log batch: %v", err)
	}

	// 由於我們沒有真正的Kafka生產者，這裡只測試函數是否執行成功
}

func TestAnonymizeLogEntry(t *testing.T) {
	// 創建Ingestor服務
	service := NewIngestorService()

	// 測試日誌條目
	logEntry := models.LogEntry{
		Service:   "test-service",
		Level:     "INFO",
		Message:   "User email: user@example.com and phone: 123-456-7890",
		Timestamp: "2023-07-15T12:00:00Z",
		UserID:    "test-user-123",
		Metadata: map[string]interface{}{
			"details": "Credit card: 4111-1111-1111-1111",
		},
	}

	// 匿名化日誌條目
	anonymized := service.anonymizeLogEntry(logEntry)

	// 檢查PII是否被匿名化
	if anonymized.Message == logEntry.Message {
		t.Error("PII was not anonymized in the message")
	}

	// 檢查是否包含[EMAIL]標記
	if !contains(anonymized.Message, "[EMAIL]") {
		t.Error("Email was not anonymized properly")
	}

	// 檢查元數據是否被匿名化
	if details, ok := anonymized.Metadata["details"].(string); ok {
		if details == "Credit card: 4111-1111-1111-1111" {
			t.Error("PII was not anonymized in metadata")
		}
		if !contains(details, "[CREDIT_CARD]") {
			t.Error("Credit card was not anonymized properly in metadata")
		}
	} else {
		t.Error("Metadata details field is missing or not a string")
	}

	// 檢查其他字段是否保持不變
	if anonymized.Service != logEntry.Service {
		t.Errorf("Service field was modified: expected %s, got %s", logEntry.Service, anonymized.Service)
	}
	if anonymized.Level != logEntry.Level {
		t.Errorf("Level field was modified: expected %s, got %s", logEntry.Level, anonymized.Level)
	}
	if anonymized.UserID != logEntry.UserID {
		t.Errorf("UserID field was modified: expected %s, got %s", logEntry.UserID, anonymized.UserID)
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