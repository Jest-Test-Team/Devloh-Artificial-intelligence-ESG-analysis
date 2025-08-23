package messaging

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaConsumerService handles Kafka message consumption
type KafkaConsumerService struct {
	reader *kafka.Reader
	handler func([]byte) error
}

// NewKafkaConsumerService creates a new Kafka consumer service
func NewKafkaConsumerService(brokers []string, topic string, groupID string, handler func([]byte) error) *KafkaConsumerService {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumerService{
		reader:  reader,
		handler: handler,
	}
}

// Start starts consuming messages from Kafka
func (k *KafkaConsumerService) Start(ctx context.Context) error {
	log.Printf("Starting Kafka consumer for topic: %s", k.reader.Config().Topic)

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer context cancelled, stopping...")
			return ctx.Err()
		default:
			m, err := k.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				time.Sleep(time.Second) // Wait before retrying
				continue
			}

			log.Printf("Received message from topic %s, partition %d, offset %d", 
				m.Topic, m.Partition, m.Offset)

			// Process the message
			if err := k.handler(m.Value); err != nil {
				log.Printf("Error processing message: %v", err)
				// Continue processing other messages
			}
		}
	}
}

// Stop stops the Kafka consumer
func (k *KafkaConsumerService) Stop() error {
	log.Println("Stopping Kafka consumer...")
	return k.reader.Close()
}

// GetStats returns consumer statistics
func (k *KafkaConsumerService) GetStats() map[string]interface{} {
	stats := k.reader.Stats()
	return map[string]interface{}{
		"messages":   stats.Messages,
		"bytes":      stats.Bytes,
		"errors":     stats.Errors,
		"rebalances": stats.Rebalances,
	}
}
