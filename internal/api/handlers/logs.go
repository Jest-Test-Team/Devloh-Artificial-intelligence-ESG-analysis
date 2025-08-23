package handlers

import (
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
		http.Error(w, "ProjectID is required", http.StatusBadRequest)
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