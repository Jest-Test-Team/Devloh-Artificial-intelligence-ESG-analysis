package processor

import (
	"devlog-ai/internal/api/models"
	"devlog-ai/internal/messaging"
	"encoding/json"
	"fmt"
	"log"
)

// Service represents the processor service
type Service struct {
	kafkaConsumer *messaging.KafkaConsumer
	stopCh        chan struct{}
	isRunning     bool
}

// NewService creates a new processor service
func NewService(kafkaConsumer *messaging.KafkaConsumer) *Service {
	return &Service{
		kafkaConsumer: kafkaConsumer,
		stopCh:        make(chan struct{}),
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
		if err := s.processLogEntry(entry); err != nil {
			log.Printf("Error processing log entry: %v", err)
			// Continue processing other entries even if one fails
		}
	}

	return nil
}

// processLogEntry processes a single log entry
func (s *Service) processLogEntry(entry models.LogEntry) error {
	// In a real implementation, we would:
	// 1. Enrich the log with additional metadata
	// 2. Apply rules or ML models for anomaly detection
	// 3. Store the processed log in a database or vector store
	
	// For now, we'll just log that we've processed it
	log.Printf("Processed log: [%s] %s: %s (source: %s, service: %s)",
		entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
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
