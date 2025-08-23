package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"devlog-ai/internal/api/models"
)

// KafkaProducer interface for producing messages to Kafka
type KafkaProducer interface {
	ProduceMessage(topic string, key string, value []byte) error
	SendLogBatch(logs []models.LogEntry) error
	Close() error
}

// IngestLogHandler handles log ingestion requests
func IngestLogHandler(producer KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var logEntry models.LogEntry
		if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set user ID from context if available
		if userID := getUserIDFromContext(r.Context()); userID != "" {
			logEntry.UserID = userID
		}

		// Set timestamps
		now := time.Now()
		logEntry.CreatedAt = now
		logEntry.UpdatedAt = now

		// Generate ID if not provided
		if logEntry.ID == "" {
			logEntry.ID = generateLogID()
		}

		// Send to Kafka
		if err := producer.ProduceMessage("logs", logEntry.ID, mustMarshal(logEntry)); err != nil {
			http.Error(w, "Failed to process log", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Log ingested successfully",
			"id":      logEntry.ID,
		})
	}
}

// getUserIDFromContext extracts user ID from context
func getUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value("user_id").(string); ok {
		return userID
	}
	return ""
}

// generateLogID generates a unique log ID
func generateLogID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// mustMarshal marshals JSON or panics (for simplicity in this example)
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// IngestLogBatchHandler handles batch log ingestion requests
func IngestLogBatchHandler(producer KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var logs []models.LogEntry
		if err := json.NewDecoder(r.Body).Decode(&logs); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Set user ID from context if available
		if userID := getUserIDFromContext(r.Context()); userID != "" {
			for i := range logs {
				logs[i].UserID = userID
			}
		}

		// Set timestamps and generate IDs
		now := time.Now()
		for i := range logs {
			if logs[i].ID == "" {
				logs[i].ID = generateLogID()
			}
			logs[i].CreatedAt = now
			logs[i].UpdatedAt = now
		}

		// Send to Kafka
		if err := producer.SendLogBatch(logs); err != nil {
			http.Error(w, "Failed to process log batch", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Log batch ingested successfully",
		})
	}
}
