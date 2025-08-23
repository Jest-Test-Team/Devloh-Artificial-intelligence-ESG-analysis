package messaging

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// 測試用的消息處理器
type testMessageHandler struct {
	messages []string
	mu       sync.Mutex
}

func newTestMessageHandler() *testMessageHandler {
	return &testMessageHandler{
		messages: make([]string, 0),
	}
}

func (h *testMessageHandler) HandleMessage(message []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if string(message) == "error" {
		return errors.New("test error")
	}
	
	h.messages = append(h.messages, string(message))
	return nil
}

func (h *testMessageHandler) GetMessages() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	result := make([]string, len(h.messages))
	copy(result, h.messages)
	return result
}

func TestMockKafkaConsumer(t *testing.T) {
	handler := newTestMessageHandler()
	consumer := NewMockKafkaConsumer(handler)
	
	// 測試啟動前添加消息（應該失敗）
	err := consumer.AddMessage([]byte("test-message-1"))
	if err == nil {
		t.Error("Expected error when adding message to stopped consumer")
	}
	
	// 測試啟動
	err = consumer.Start()
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}
	
	// 測試添加消息
	err = consumer.AddMessage([]byte("test-message-2"))
	if err != nil {
		t.Fatalf("Failed to add message: %v", err)
	}
	
	// 檢查消息是否被處理
	messages := handler.GetMessages()
	if len(messages) != 1 || messages[0] != "test-message-2" {
		t.Errorf("Expected message to be processed, got: %v", messages)
	}
	
	// 測試添加錯誤消息
	err = consumer.AddMessage([]byte("error"))
	if err == nil {
		t.Error("Expected error when handling error message")
	}
	
	// 測試停止
	err = consumer.Stop()
	if err != nil {
		t.Fatalf("Failed to stop consumer: %v", err)
	}
	
	// 測試停止後添加消息（應該失敗）
	err = consumer.AddMessage([]byte("test-message-3"))
	if err == nil {
		t.Error("Expected error when adding message to stopped consumer")
	}
}

func TestNewKafkaConsumer(t *testing.T) {
	handler := newTestMessageHandler()
	
	// 測試創建消費者
	consumer, err := NewKafkaConsumer([]string{"test-topic"}, "test-group", handler)
	if err != nil {
		t.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	
	// 測試啟動
	err = consumer.Start()
	if err != nil {
		t.Fatalf("Failed to start consumer: %v", err)
	}
	
	// 等待一小段時間，確保消費者啟動
	time.Sleep(100 * time.Millisecond)
	
	// 測試停止
	err = consumer.Stop()
	if err != nil {
		t.Fatalf("Failed to stop consumer: %v", err)
	}
}