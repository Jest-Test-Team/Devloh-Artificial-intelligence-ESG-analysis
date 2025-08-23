package messaging

import (
	"context"
	"crypto/tls"
	"devlog-ai/internal/api/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// RealKafkaProducer is a real implementation of the KafkaProducer interface
type RealKafkaProducer struct {
	writer *kafka.Writer
	config KafkaConfig
}

// RealKafkaConsumer is a real implementation of the KafkaConsumer interface
type RealKafkaConsumer struct {
	reader  *kafka.Reader
	config  KafkaConfig
	handler func([]byte) error
	stopCh  chan struct{}
}

// KafkaConfig holds configuration for Kafka
type KafkaConfig struct {
	Brokers  []string
	Username string
	Password string
	Topic    string
	GroupID  string // For consumer
}

// NewRealKafkaProducer creates a new real Kafka producer
func NewRealKafkaProducer(config KafkaConfig) (*RealKafkaProducer, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Configure TLS and SASL if username and password are provided
	if config.Username != "" && config.Password != "" {
		dialer.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		dialer.SASLMechanism = kafka.PLAIN(config.Username, config.Password)
	}

	// Create a writer
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	}

	return &RealKafkaProducer{
		writer: writer,
		config: config,
	}, nil
}

// SendLogBatch sends a batch of logs to Kafka
func (p *RealKafkaProducer) SendLogBatch(batch models.LogBatch) error {
	// Convert batch to JSON
	jsonData, err := json.Marshal(batch)
	if err != nil {
		return fmt.Errorf("error marshaling log batch: %w", err)
	}

	// Send to Kafka
	msg := kafka.Message{
		Key:   []byte(batch.ProjectID),
		Value: jsonData,
		Time:  time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	log.Printf("Sent %d logs to Kafka topic %s for project %s",
		len(batch.Logs), p.config.Topic, batch.ProjectID)

	return nil
}

// Close closes the Kafka producer
func (p *RealKafkaProducer) Close() error {
	return p.writer.Close()
}

// NewRealKafkaConsumer creates a new real Kafka consumer
func NewRealKafkaConsumer(config KafkaConfig) (*RealKafkaConsumer, error) {
	if len(config.Brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}
	if config.GroupID == "" {
		return nil, fmt.Errorf("no group ID provided")
	}

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	// Configure TLS and SASL if username and password are provided
	if config.Username != "" && config.Password != "" {
		dialer.TLS = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
		dialer.SASLMechanism = kafka.PLAIN(config.Username, config.Password)
	}

	// Create a reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        config.Brokers,
		GroupID:        config.GroupID,
		Topic:          config.Topic,
		MinBytes:       10e3,    // 10KB
		MaxBytes:       10e6,    // 10MB
		MaxWait:        1 * time.Second,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 1 * time.Second,
		Dialer:         dialer,
	})

	return &RealKafkaConsumer{
		reader: reader,
		config: config,
		stopCh: make(chan struct{}),
	}, nil
}

// Start starts consuming messages
func (c *RealKafkaConsumer) Start(topic string, handler func([]byte) error) error {
	if topic != c.config.Topic {
		return fmt.Errorf("topic mismatch: configured %s, requested %s", c.config.Topic, topic)
	}

	c.handler = handler

	// Start consuming in a goroutine
	go func() {
		for {
			select {
			case <-c.stopCh:
				return
			default:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				msg, err := c.reader.ReadMessage(ctx)
				cancel()

				if err != nil {
					if err == context.DeadlineExceeded {
						// This is normal, just continue
						continue
					}
					log.Printf("Error reading message: %v", err)
					continue
				}

				// Handle the message
				if err := c.handler(msg.Value); err != nil {
					log.Printf("Error handling message: %v", err)
				}
			}
		}
	}()

	log.Printf("Started Kafka consumer for topic: %s", c.config.Topic)
	return nil
}

// Stop stops consuming messages
func (c *RealKafkaConsumer) Stop() {
	close(c.stopCh)
	c.reader.Close()
	log.Println("Kafka consumer stopped")
}

// GetKafkaConfigFromEnv gets Kafka configuration from environment variables
func GetKafkaConfigFromEnv() (KafkaConfig, error) {
	config := KafkaConfig{
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: os.Getenv("KAFKA_GROUP_ID"),
	}

	// Get the broker list
	brokersStr := os.Getenv("KAFKA_BROKERS")
	if brokersStr == "" {
		return config, fmt.Errorf("KAFKA_BROKERS environment variable not set")
	}

	// Parse the broker list
	brokers := ParseStringList(brokersStr)
	if len(brokers) == 0 {
		return config, fmt.Errorf("no Kafka brokers specified")
	}
	config.Brokers = brokers

	// Get username and password
	config.Username = os.Getenv("KAFKA_USERNAME")
	config.Password = os.Getenv("KAFKA_PASSWORD")

	// Set defaults if needed
	if config.Topic == "" {
		config.Topic = LogsTopic
	}
	if config.GroupID == "" {
		config.GroupID = "devlogai-processor"
	}

	return config, nil
}
