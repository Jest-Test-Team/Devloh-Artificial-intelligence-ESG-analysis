package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"devlog-ai/internal/api/models"
)

// Topics
const (
	LogsTopic = "logs"
)

// KafkaProducer handles sending messages to Kafka
type KafkaProducer struct {
	// This would normally hold the Kafka producer client
	// For now, we just simulate it for development purposes
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer() (*KafkaProducer, error) {
	// TODO: Initialize real Kafka producer when we have Kafka setup
	// For now, just return a mock producer
	return &KafkaProducer{}, nil
}

// SendLogBatch sends a batch of logs to Kafka
func (p *KafkaProducer) SendLogBatch(batch models.LogBatch) error {
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
func (p *KafkaProducer) Close() error {
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