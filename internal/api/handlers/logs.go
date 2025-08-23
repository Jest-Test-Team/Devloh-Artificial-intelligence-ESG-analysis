package handlers

import (
	"devlog-ai/internal/api/middleware"
	"devlog-ai/internal/api/models"
	"devlog-ai/internal/messaging"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var kafkaProducer *messaging.KafkaProducer

// InitLogHandlers initializes dependencies for log handlers
func InitLogHandlers() error {
	// Initialize Kafka producer
	var err error
	kafkaProducer, err = messaging.NewKafkaProducer()
	if err != nil {
		return err
	}
	
	return nil
}

// IngestLogsHandler receives log batches and forwards them for processing
func IngestLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the log batch
	var logBatch models.LogBatch
	if err := json.Unmarshal(body, &logBatch); err != nil {
		http.Error(w, "Invalid log format", http.StatusBadRequest)
		return
	}

	// Validate the log batch
	if logBatch.ProjectID == "" {
		// If no project ID was provided, use the one from the user
		logBatch.ProjectID = user.ProjectID
	}
	
	// Ensure the user can only send logs to their own project
	if logBatch.ProjectID != user.ProjectID {
		http.Error(w, "Unauthorized: can only send logs to your own project", http.StatusForbidden)
		return
	}

	if len(logBatch.Logs) == 0 {
		http.Error(w, "Log batch cannot be empty", http.StatusBadRequest)
		return
	}

	// TODO: Anonymize PII data as per risk mitigation guidelines

	// Send logs to Kafka for processing
	if kafkaProducer != nil {
		if err := kafkaProducer.SendLogBatch(logBatch); err != nil {
			log.Printf("Error sending logs to Kafka: %v", err)
			http.Error(w, "Error processing logs", http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("Warning: Kafka producer not initialized, logs will not be processed")
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}