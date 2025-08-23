package ingestor

import (
	"devlog-ai/internal/api/models"
	"devlog-ai/internal/messaging"
	"devlog-ai/internal/privacy"
	"log"
	"time"
)

// Service represents the ingestor service
type Service struct {
	kafkaProducer messaging.KafkaProducer // 改為介面類型，而非指針到介面
	stopCh        chan struct{}
	isRunning     bool
	enablePII     bool // Flag to enable/disable PII anonymization
}

// ServiceConfig holds configuration options for the ingestor service
type ServiceConfig struct {
	EnablePII bool // Whether to enable PII anonymization
}

// NewService creates a new ingestor service
func NewService(kafkaProducer messaging.KafkaProducer, config ServiceConfig) *Service {
	return &Service{
		kafkaProducer: kafkaProducer,
		stopCh:        make(chan struct{}),
		enablePII:     config.EnablePII,
	}
}

// Start starts the service
func (s *Service) Start() error {
	if s.isRunning {
		return nil
	}

	s.isRunning = true
	go s.run()
	log.Println("Ingestor service started")
	return nil
}

// Stop stops the service
func (s *Service) Stop() {
	if !s.isRunning {
		return
	}

	close(s.stopCh)
	s.isRunning = false
	log.Println("Ingestor service stopped")
}

// run is the main processing loop
func (s *Service) run() {
	// In a real implementation, this would listen for logs from various sources
	// For now, we'll just simulate periodic activity
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Simulated activity
			log.Println("Ingestor is running, ready to process logs")
		case <-s.stopCh:
			return
		}
	}
}

// ProcessLogBatch processes a batch of logs
func (s *Service) ProcessLogBatch(batch models.LogBatch) error {
	// Perform any preprocessing on logs
	
	// Ensure all logs have timestamps
	for i := range batch.Logs {
		if batch.Logs[i].Timestamp.IsZero() {
			batch.Logs[i].Timestamp = time.Now()
		}

		// Apply PII anonymization if enabled
		if s.enablePII {
			batch.Logs[i].Message = privacy.AnonymizePII(batch.Logs[i].Message)
			
			// Also anonymize metadata if present
			if batch.Logs[i].Metadata != nil {
				for key, value := range batch.Logs[i].Metadata {
					batch.Logs[i].Metadata[key] = privacy.AnonymizePII(value)
				}
			}
		}
	}

	// Send to Kafka
	if err := s.kafkaProducer.SendLogBatch(batch); err != nil {
		return err
	}

	return nil
}