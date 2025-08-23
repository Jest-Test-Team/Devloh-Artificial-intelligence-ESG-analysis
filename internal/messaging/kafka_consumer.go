package messaging

import (
	"devlog-ai/internal/api/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// KafkaConsumer handles consuming messages from Kafka
type KafkaConsumer struct {
	// This would normally hold the Kafka consumer client
	// For now, we just simulate it for development purposes
	stopCh chan struct{}
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer() (*KafkaConsumer, error) {
	// TODO: Initialize real Kafka consumer when we have Kafka setup
	// For now, just return a mock consumer
	return &KafkaConsumer{
		stopCh: make(chan struct{}),
	}, nil
}

// Start starts consuming messages
func (c *KafkaConsumer) Start(topic string, handler func([]byte) error) error {
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

	log.Printf("Started Kafka consumer for topic: %s", topic)
	return nil
}

// Stop stops consuming messages
func (c *KafkaConsumer) Stop() {
	close(c.stopCh)
	log.Println("Kafka consumer stopped")
}
