package processor

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"devlog-ai/internal/api/models"
	"devlog-ai/internal/embedding"
	"devlog-ai/internal/vectordb"
)

// 測試用的消息處理器
type testMessageHandler struct {
	messages []models.LogEntry
	callback func(models.LogEntry)
}

func (h *testMessageHandler) HandleMessage(msg []byte) error {
	var logEntry models.LogEntry
	if err := json.Unmarshal(msg, &logEntry); err != nil {
		return err
	}
	h.messages = append(h.messages, logEntry)
	if h.callback != nil {
		h.callback(logEntry)
	}
	return nil
}

// 測試用的模擬Kafka消費者
type testKafkaConsumer struct {
	handler MessageHandler
	started bool
	stopped bool
}

func (c *testKafkaConsumer) Start(ctx context.Context) error {
	c.started = true
	return nil
}

func (c *testKafkaConsumer) Stop() error {
	c.stopped = true
	return nil
}

func (c *testKafkaConsumer) ConsumeMessages(messages [][]byte) {
	for _, msg := range messages {
		c.handler.HandleMessage(msg)
	}
}

func TestLogProcessorService_Start(t *testing.T) {
	// 創建模擬嵌入客戶端和向量數據庫客戶端
	mockEmbedding := embedding.NewMockEmbeddingClient()
	mockVectorDB := vectordb.NewMockVectorDBClient()

	// 創建測試消息處理器
	handler := &testMessageHandler{
		messages: make([]models.LogEntry, 0),
	}

	// 創建測試Kafka消費者
	consumer := &testKafkaConsumer{
		handler: handler,
	}

	// 創建處理器服務
	service := NewLogProcessorService(consumer, mockEmbedding, mockVectorDB)

	// 啟動服務
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := service.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start processor service: %v", err)
	}

	// 檢查消費者是否已啟動
	if !consumer.started {
		t.Error("Kafka consumer was not started")
	}

	// 停止服務
	err = service.Stop()
	if err != nil {
		t.Fatalf("Failed to stop processor service: %v", err)
	}

	// 檢查消費者是否已停止
	if !consumer.stopped {
		t.Error("Kafka consumer was not stopped")
	}
}

func TestLogProcessorService_ProcessLogEntry(t *testing.T) {
	// 創建模擬嵌入客戶端和向量數據庫客戶端
	mockEmbedding := embedding.NewMockEmbeddingClient()
	mockVectorDB := vectordb.NewMockVectorDBClient()

	// 創建測試消息處理器
	var processedLog models.LogEntry
	handler := &testMessageHandler{
		messages: make([]models.LogEntry, 0),
		callback: func(log models.LogEntry) {
			processedLog = log
		},
	}

	// 創建測試Kafka消費者
	consumer := &testKafkaConsumer{
		handler: handler,
	}

	// 創建處理器服務
	service := NewLogProcessorService(consumer, mockEmbedding, mockVectorDB)

	// 創建測試日誌
	logEntry := models.LogEntry{
		ID:        "test-log-123",
		Service:   "test-service",
		Level:     "INFO",
		Message:   "This is a test log message",
		Timestamp: time.Now().Format(time.RFC3339),
		UserID:    "test-user-456",
		Metadata: map[string]interface{}{
			"request_id": "req-789",
		},
	}

	// 序列化日誌
	logBytes, err := json.Marshal(logEntry)
	if err != nil {
		t.Fatalf("Failed to marshal log entry: %v", err)
	}

	// 模擬消費消息
	consumer.ConsumeMessages([][]byte{logBytes})

	// 檢查日誌是否被處理
	if processedLog.ID != logEntry.ID {
		t.Errorf("Log was not processed correctly: expected ID %s, got %s", logEntry.ID, processedLog.ID)
	}

	// 檢查向量數據庫是否有插入
	results, err := mockVectorDB.Query([]float32{0.1, 0.2, 0.3}, 1)
	if err != nil {
		t.Fatalf("Failed to query vector database: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result from vector database, got %d", len(results))
	}
}
