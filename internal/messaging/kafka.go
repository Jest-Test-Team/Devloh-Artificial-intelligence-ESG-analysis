package messaging

import (
	"devlog-ai/internal/api/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// Topics
const (
	LogsTopic = "logs"
)

// KafkaProducer handles sending messages to Kafka
type KafkaProducer interface {
	// SendLogBatch sends a batch of logs to Kafka
	SendLogBatch(batch models.LogBatch) error
	
	// Close closes the Kafka producer
	Close() error
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer() (KafkaProducer, error) {
	// Check if we should use a real Kafka producer
	useReal := os.Getenv("USE_REAL_KAFKA")
	if strings.ToLower(useReal) == "true" {
		// Get Kafka configuration from environment variables
		config, err := GetKafkaConfigFromEnv()
		if err != nil {
			return nil, fmt.Errorf("failed to get Kafka configuration: %w", err)
		}
		
		// Create a real Kafka producer
		producer, err := NewRealKafkaProducer(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create real Kafka producer: %w", err)
		}
		
		log.Println("Using real Kafka producer")
		return producer, nil
	}
	
	// Use a mock Kafka producer
	log.Println("Using mock Kafka producer")
	return NewMockKafkaProducer(), nil
}

// MockKafkaProducer is a mock implementation of the KafkaProducer interface
type MockKafkaProducer struct {
	// This would normally hold the Kafka producer client
	// For now, we just simulate it for development purposes
}

// NewMockKafkaProducer creates a new mock Kafka producer
func NewMockKafkaProducer() *MockKafkaProducer {
	return &MockKafkaProducer{}
}

// SendLogBatch sends a batch of logs to Kafka
func (p *MockKafkaProducer) SendLogBatch(batch models.LogBatch) error {
	// Convert batch to JSON
	jsonData, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("error marshaling log batch: %w", err)
	}

	// In a real implementation, we would send this to Kafka
	// For now, just log it
	log.Printf("Would send to Kafka topic '%s': %d bytes of log data for project %s",
		LogsTopic, len(jsonData), batch.ProjectID)

	// Check for DEBUG environment variable to print more details
	if os.Getenv("DEBUG") == "true" {
		log.Printf("Log data (first 200 chars): %s", string(jsonData)[:min(200, len(jsonData))])
	}

	return nil
}

// Close closes the Kafka producer
func (p *MockKafkaProducer) Close() error {
	// TODO: Close the real Kafka producer when we have it
	return nil
}

// Helper function for string truncation
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ParseStringList parses a comma-separated string into a slice of strings
func ParseStringList(s string) []string {
	if s == "" {
		return nil
	}
	
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	
	return result
}