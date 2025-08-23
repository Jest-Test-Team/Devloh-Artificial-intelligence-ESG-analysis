package messaging

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"devlog-ai/internal/api/models"
)

// 測試用的消息處理器
type testMessageHandler struct {
	messages []models.LogEntry
}

func (h *testMessageHandler) HandleMessage(msg []byte) error {
	var logEntry models.LogEntry
	if err := json.Unmarshal(msg, &logEntry); err != nil {
		return err
	}
	h.messages = append(h.messages, logEntry)
	return nil
}

func TestMockKafkaConsumer(t *testing.T) {
	// 創建測試消息處理器
	handler := &testMessageHandler{
		messages: make([]models.LogEntry, 0),
	}

	// 創建模擬Kafka消費者
	consumer := NewMockKafkaConsumer("test-topic", "test-group", handler)

	// 啟動消費者
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 啟動消費者
	err := consumer.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}

	// 模擬發送消息
	logEntry := models.LogEntry{
		Service:   "test-service",
		Level:     "INFO",
		Message:   "Test message",
		Timestamp: "2023-07-15T12:00:00Z",
	}
	logBytes, _ := json.Marshal(logEntry)

	// 使用ConsumeMessages方法模擬接收消息
	consumer.(*MockKafkaConsumer).ConsumeMessages([][]byte{logBytes})

	// 等待消息處理
	time.Sleep(100 * time.Millisecond)

	// 檢查消息是否被處理
	if len(handler.messages) != 1 {
		t.Errorf("Expected 1 message to be processed, got %d", len(handler.messages))
	}

	// 檢查消息內容
	if handler.messages[0].Service != "test-service" {
		t.Errorf("Expected service 'test-service', got '%s'", handler.messages[0].Service)
	}

	// 停止消費者
	err = consumer.Stop()
	if err != nil {
		t.Fatalf("Failed to stop consumer: %v", err)
	}
}

func TestNewKafkaConsumer(t *testing.T) {
	// 創建測試消息處理器
	handler := &testMessageHandler{
		messages: make([]models.LogEntry, 0),
	}

	// 測試創建Kafka消費者
	consumer, err := NewKafkaConsumer("test-topic", "test-group", handler)
	if err != nil {
		t.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// 檢查是否是模擬消費者（因為沒有設置環境變量）
	if _, ok := consumer.(*MockKafkaConsumer); !ok {
		t.Error("Expected MockKafkaConsumer when environment variables are not set")
	}
}
