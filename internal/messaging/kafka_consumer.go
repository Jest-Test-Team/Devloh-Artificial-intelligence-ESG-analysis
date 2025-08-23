package messaging

import (
	"devlog-ai/internal/api/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// KafkaConsumer handles consuming messages from Kafka
type KafkaConsumer interface {
	// Start starts consuming messages
	Start(topic string, handler func([]byte) error) error
	
	// Stop stops consuming messages
	Stop()
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer() (KafkaConsumer, error) {
	// Check if we should use a real Kafka consumer
	useReal := os.Getenv("USE_REAL_KAFKA")
	if strings.ToLower(useReal) == "true" {
		// Get Kafka configuration from environment variables
		config, err := GetKafkaConfigFromEnv()
		if err != nil {
			return nil, fmt.Errorf("failed to get Kafka configuration: %w", err)
		}
		
		// Create a real Kafka consumer
		consumer, err := NewRealKafkaConsumer(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create real Kafka consumer: %w", err)
		}
		
		log.Println("Using real Kafka consumer")
		return consumer, nil
	}
	
	// Use a mock Kafka consumer
	log.Println("Using mock Kafka consumer")
	return NewMockKafkaConsumer(), nil
}

// MockKafkaConsumer is a mock implementation of the KafkaConsumer interface
type MockKafkaConsumer struct {
	// This would normally hold the Kafka consumer client
	// For now, we just simulate it for development purposes
	stopCh  chan struct{}
	handler func([]byte) error
}

// NewMockKafkaConsumer creates a new mock Kafka consumer
func NewMockKafkaConsumer() *MockKafkaConsumer {
	return &MockKafkaConsumer{
		stopCh: make(chan struct{}),
	}
}

// Start starts consuming messages
func (c *MockKafkaConsumer) Start(topic string, handler func([]byte) error) error {
	c.handler = handler

	// In a real implementation, this would start consuming from Kafka
	// For now, we'll simulate periodic message reception
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Simulate receiving a message
				if topic == LogsTopic {
					// Create a simulated log batch
					batch := models.LogBatch{
						ProjectID: "sim-project",
						Logs: []models.LogEntry{
							{
								Timestamp:   time.Now(),
								Level:       "INFO",
								Message:     "Simulated log entry for testing",
								Source:      "simulator",
								ServiceName: "test-service",
							},
						},
					}

					// Convert to JSON
					jsonData, err := json.Marshal(batch)
					if err != nil {
						log.Printf("Error creating simulated message: %v", err)
						continue
					}

					// Call the handler
					if err := handler(jsonData); err != nil {
						log.Printf("Error handling simulated message: %v", err)
					} else {
						log.Printf("Successfully processed simulated message on topic: %s", topic)
					}
				}
			case <-c.stopCh:
				return
			}
		}
	}()

	log.Printf("Started mock Kafka consumer for topic: %s", topic)
	return nil
}

// Stop stops consuming messages
func (c *MockKafkaConsumer) Stop() {
	close(c.stopCh)
	log.Println("Mock Kafka consumer stopped")
}