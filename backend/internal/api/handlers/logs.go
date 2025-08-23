package handlers

import (
	"encoding/json"
	"net/http"

	"devlog-ai/internal/api/middleware"
	"devlog-ai/internal/api/models"
)

// KafkaProducer 是一個 Kafka 生產者的接口
type KafkaProducer interface {
	ProduceMessage(topic string, key string, value []byte) error
	SendLogBatch(logs []models.LogEntry) error
	Close() error
}

// IngestLogHandler 處理單條日誌的攝取
func IngestLogHandler(producer KafkaProducer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析請求體
		var logEntry models.LogEntry
		if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 從上下文中獲取用戶ID
		userID := middleware.GetUserID(r.Context())
		logEntry.UserID = userID

		// 將日誌序列化為JSON
		logBytes, err := json.Marshal(logEntry)
		if err != nil {
			http.Error(w, "Failed to marshal log entry", http.StatusInternalServerError)
			return
		}

		// 發送到Kafka
		if err := producer.ProduceMessage("logs", logEntry.Service, logBytes); err != nil {
			http.Error(w, "Failed to produce message", http.StatusInternalServerError)
			return
		}

		// 返回成功響應
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}

// IngestLogBatchHandler 處理批量日誌的攝取
func IngestLogBatchHandler(producer KafkaProducer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 解析請求體
		var logEntries []models.LogEntry
		if err := json.NewDecoder(r.Body).Decode(&logEntries); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 從上下文中獲取用戶ID
		userID := middleware.GetUserID(r.Context())

		// 設置每條日誌的用戶ID
		for i := range logEntries {
			logEntries[i].UserID = userID
		}

		// 發送到Kafka
		if err := producer.SendLogBatch(logEntries); err != nil {
			http.Error(w, "Failed to send log batch", http.StatusInternalServerError)
			return
		}

		// 返回成功響應
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}