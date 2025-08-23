package processor

import (
	"devlog-ai/internal/api/models"
	"devlog-ai/internal/embedding"
	"devlog-ai/internal/messaging"
	"devlog-ai/internal/vectordb"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Service represents the processor service
type Service struct {
	kafkaConsumer  *messaging.KafkaConsumer
	embeddingClient embedding.Client
	vectorDB       vectordb.Client
	stopCh         chan struct{}
	isRunning      bool
}

// ServiceConfig holds configuration options for the processor service
type ServiceConfig struct {
	EnableVectorDB bool // Whether to enable vector database integration
}

// NewService creates a new processor service
func NewService(kafkaConsumer *messaging.KafkaConsumer, embeddingClient embedding.Client, 
	vectorDB vectordb.Client, config ServiceConfig) *Service {
	return &Service{
		kafkaConsumer:   kafkaConsumer,
		embeddingClient: embeddingClient,
		vectorDB:        vectorDB,
		stopCh:          make(chan struct{}),
	}
}

// Start starts the service
func (s *Service) Start() error {
	if s.isRunning {
		return nil
	}

	// Start consuming log messages from Kafka
	if err := s.kafkaConsumer.Start(messaging.LogsTopic, s.processLogMessage); err != nil {
		return fmt.Errorf("failed to start Kafka consumer: %w", err)
	}

	s.isRunning = true
	log.Println("Processor service started")
	return nil
}

// Stop stops the service
func (s *Service) Stop() {
	if !s.isRunning {
		return
	}

	s.kafkaConsumer.Stop()
	close(s.stopCh)
	s.isRunning = false
	log.Println("Processor service stopped")
}

// processLogMessage handles a log message from Kafka
func (s *Service) processLogMessage(message []byte) error {
	// Unmarshal the message
	var logBatch models.LogBatch
	if err := json.Unmarshal(message, &logBatch); err != nil {
		return fmt.Errorf("failed to unmarshal log batch: %w", err)
	}

	// Process each log entry
	for _, entry := range logBatch.Logs {
		if err := s.processLogEntry(entry, logBatch.ProjectID); err != nil {
			log.Printf("Error processing log entry: %v", err)
			// Continue processing other entries even if one fails
		}
	}

	return nil
}

// processLogEntry processes a single log entry
func (s *Service) processLogEntry(entry models.LogEntry, projectID string) error {
	// Combine the log fields into a single text for embedding
	logText := fmt.Sprintf("[%s] [%s] %s (Source: %s, Service: %s)", 
		entry.Timestamp.Format(time.RFC3339), entry.Level, entry.Message, 
		entry.Source, entry.ServiceName)

	// Store the embedding in the vector database if available
	if s.embeddingClient != nil && s.vectorDB != nil {
		// Generate embedding
		vector, err := s.embeddingClient.EmbedText(logText)
		if err != nil {
			log.Printf("Error generating embedding: %v", err)
			// Continue processing even if embedding fails
		} else {
			// Create metadata for the vector database
			metadata := map[string]interface{}{
				"project_id":   projectID,
				"timestamp":    entry.Timestamp.Format(time.RFC3339),
				"level":        entry.Level,
				"message":      entry.Message,
				"source":       entry.Source,
				"service_name": entry.ServiceName,
			}

			// Add any additional metadata from the log entry
			if entry.Metadata != nil {
				for key, value := range entry.Metadata {
					metadata[key] = value
				}
			}

			// Generate a unique ID for the entry
			id := fmt.Sprintf("%s-%s", projectID, entry.ID)
			if entry.ID == "" {
				id = fmt.Sprintf("%s-%d", projectID, time.Now().UnixNano())
			}

			// Store the embedding in the vector database
			if err := s.vectorDB.UpsertVector(id, vector, metadata); err != nil {
				log.Printf("Error storing embedding: %v", err)
				// Continue processing even if storage fails
			} else {
				log.Printf("Successfully stored embedding for log entry %s", id)
			}
		}
	}

	// Perform anomaly detection or other processing here

	log.Printf("Processed log: [%s] %s: %s (source: %s, service: %s)",
		entry.Timestamp.Format(time.RFC3339),
		entry.Level,
		entry.Message,
		entry.Source,
		entry.ServiceName,
	)

	// If this is an ERROR level log, we might want to trigger an alert
	if entry.Level == "ERROR" {
		log.Printf("Alert: Error detected in service %s: %s", entry.ServiceName, entry.Message)
		// TODO: Implement actual alerting system
	}

	return nil
}