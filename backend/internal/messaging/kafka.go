package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer produces messages to Kafka
type KafkaProducer struct {
	writer *kafka.Writer
}

// KafkaConsumer consumes messages from Kafka
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
	}
}

// NewKafkaConsumer creates a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic string, groupID string) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaConsumer{
		reader: reader,
	}
}

// ProduceMessage produces a message to Kafka
func (k *KafkaProducer) ProduceMessage(topic string, key string, value []byte) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	return k.writer.WriteMessages(context.Background(), msg)
}

// SendLogBatch sends a batch of log entries to Kafka
func (k *KafkaProducer) SendLogBatch(logs interface{}) error {
	// For now, just return success
	// TODO: Implement actual batch sending logic
	return nil
}

// Close closes the Kafka producer
func (k *KafkaProducer) Close() error {
	return k.writer.Close()
}

// ConsumeMessages consumes messages from Kafka
func (k *KafkaConsumer) ConsumeMessages(ctx context.Context, handler func([]byte) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := k.reader.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("failed to read message: %w", err)
			}

			if err := handler(m.Value); err != nil {
				// Log error but continue processing
				fmt.Printf("Error processing message: %v\n", err)
			}
		}
	}
}

// Close closes the Kafka consumer
func (k *KafkaConsumer) Close() error {
	return k.reader.Close()
}
