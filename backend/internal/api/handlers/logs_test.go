package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"devlog-ai/internal/api/middleware"
	"devlog-ai/internal/api/models"
)

// 模擬Kafka生產者
type mockKafkaProducer struct {
	messages []models.LogEntry
}

func (m *mockKafkaProducer) ProduceMessage(topic string, key string, value []byte) error {
	var logEntry models.LogEntry
	json.Unmarshal(value, &logEntry)
	m.messages = append(m.messages, logEntry)
	return nil
}

func (m *mockKafkaProducer) SendLogBatch(logs []models.LogEntry) error {
	m.messages = append(m.messages, logs...)
	return nil
}

func (m *mockKafkaProducer) Close() error {
	return nil
}

func TestIngestLogHandler(t *testing.T) {
	// 創建模擬Kafka生產者
	mockProducer := &mockKafkaProducer{
		messages: make([]models.LogEntry, 0),
	}

	// 創建日誌處理程序
	handler := IngestLogHandler(mockProducer)

	// 創建請求體
	requestBody := models.LogEntry{
		Service:   "test-service",
		Level:     "INFO",
		Message:   "這是一條測試日誌",
		Timestamp: "2023-07-15T12:00:00Z",
		Metadata: map[string]interface{}{
			"request_id": "req-123",
			"user_id":    "user-456",
		},
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// 創建一個新的請求
	req, err := http.NewRequest("POST", "/logs", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 模擬認證中間件
	ctx := req.Context()
	ctx = middleware.WithUserID(ctx, "test-user-123")
	req = req.WithContext(ctx)

	// 創建一個ResponseRecorder來記錄響應
	rr := httptest.NewRecorder()

	// 直接調用處理函數
	handler.ServeHTTP(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢查Kafka生產者是否收到消息
	if len(mockProducer.messages) != 1 {
		t.Errorf("producer did not receive message: got %v messages", len(mockProducer.messages))
	}

	// 檢查消息內容
	producedMsg := mockProducer.messages[0]
	if producedMsg.Service != "test-service" {
		t.Errorf("unexpected service: got %v want %v", producedMsg.Service, "test-service")
	}
	if producedMsg.Message != "這是一條測試日誌" {
		t.Errorf("unexpected message: got %v want %v", producedMsg.Message, "這是一條測試日誌")
	}
	if producedMsg.UserID != "test-user-123" {
		t.Errorf("user ID not set correctly: got %v want %v", producedMsg.UserID, "test-user-123")
	}
}

// 測試批量日誌處理
func TestIngestLogBatchHandler(t *testing.T) {
	// 創建模擬Kafka生產者
	mockProducer := &mockKafkaProducer{
		messages: make([]models.LogEntry, 0),
	}

	// 創建日誌處理程序
	handler := IngestLogBatchHandler(mockProducer)

	// 創建請求體
	requestBody := []models.LogEntry{
		{
			Service:   "test-service-1",
			Level:     "INFO",
			Message:   "這是第一條測試日誌",
			Timestamp: "2023-07-15T12:00:00Z",
		},
		{
			Service:   "test-service-2",
			Level:     "WARN",
			Message:   "這是第二條測試日誌",
			Timestamp: "2023-07-15T12:01:00Z",
		},
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// 創建一個新的請求
	req, err := http.NewRequest("POST", "/logs/batch", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 模擬認證中間件
	ctx := req.Context()
	ctx = middleware.WithUserID(ctx, "test-user-123")
	req = req.WithContext(ctx)

	// 創建一個ResponseRecorder來記錄響應
	rr := httptest.NewRecorder()

	// 直接調用處理函數
	handler.ServeHTTP(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢查Kafka生產者是否收到消息
	if len(mockProducer.messages) != 2 {
		t.Errorf("producer did not receive correct number of messages: got %v want %v", 
			len(mockProducer.messages), 2)
	}

	// 檢查每條消息的用戶ID
	for _, msg := range mockProducer.messages {
		if msg.UserID != "test-user-123" {
			t.Errorf("user ID not set correctly: got %v want %v", msg.UserID, "test-user-123")
		}
	}
}
