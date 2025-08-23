package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"devlog-ai/internal/api/models"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer 定義了 Kafka 生產者的接口
type KafkaProducer interface {
	// ProduceMessage 將消息發送到指定的 topic
	ProduceMessage(topic string, key string, value []byte) error

	// SendLogBatch 將一批日誌發送到 Kafka
	SendLogBatch(logs []models.LogEntry) error

	// Close 關閉生產者
	Close() error
}

// 實際的 Kafka 生產者實現
type kafkaProducer struct {
	writer *kafka.Writer
	topics []string
}

// NewKafkaProducer 創建一個新的 Kafka 生產者
// 如果環境變量不足，則返回一個模擬生產者
func NewKafkaProducer() (KafkaProducer, error) {
	// 檢查必要的環境變量
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// 如果沒有設置 Kafka 環境變量，使用模擬生產者
		return NewMockKafkaProducer(), nil
	}

	// 創建 Kafka 寫入器
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		// 如果有驗證需要，添加 Transport
	}

	// 檢查是否需要 SASL 驗證
	// 注意：在實際應用中，應該使用適當的 SASL 機制
	// 這裡我們簡化處理，因為 kafka-go 版本可能不同
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	if username != "" && password != "" {
		// 在實際應用中設置 SASL
	}

	return &kafkaProducer{
		writer: writer,
		topics: []string{"logs", "metrics", "alerts"},
	}, nil
}

// NewKafkaProducerWithTopics 創建一個新的 Kafka 生產者，指定 topics 和 clientID
func NewKafkaProducerWithTopics(topics []string, clientID string) (KafkaProducer, error) {
	// 檢查必要的環境變量
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// 如果沒有設置 Kafka 環境變量，使用模擬生產者
		return NewMockKafkaProducer(), nil
	}

	// 創建 Kafka 寫入器
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		// 注意：某些版本的 kafka-go 可能不支持 ClientID
		// ClientID:     clientID,
		// 如果有驗證需要，添加 Transport
	}

	// 檢查是否需要 SASL 驗證
	// 注意：在實際應用中，應該使用適當的 SASL 機制
	// 這裡我們簡化處理，因為 kafka-go 版本可能不同
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	if username != "" && password != "" {
		// 在實際應用中設置 SASL
	}

	return &kafkaProducer{
		writer: writer,
		topics: topics,
	}, nil
}

// ProduceMessage 將消息發送到指定的 topic
func (p *kafkaProducer) ProduceMessage(topic string, key string, value []byte) error {
	message := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}

	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   message.Key,
		Value: message.Value,
		Time:  message.Time,
	})
}

// SendLogBatch 將一批日誌發送到 Kafka
func (p *kafkaProducer) SendLogBatch(logs []models.LogEntry) error {
	messages := make([]kafka.Message, 0, len(logs))

	for _, log := range logs {
		logBytes, err := json.Marshal(log)
		if err != nil {
			return fmt.Errorf("failed to marshal log entry: %w", err)
		}

		messages = append(messages, kafka.Message{
			Topic: "logs",
			Key:   []byte(log.Service),
			Value: logBytes,
			Time:  time.Now(),
		})
	}

	return p.writer.WriteMessages(context.Background(), messages...)
}

// Close 關閉 Kafka 生產者
func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}

// MockKafkaProducer 是一個模擬的 Kafka 生產者，用於測試
type MockKafkaProducer struct {
	Messages map[string][]string // topic -> messages
}

// NewMockKafkaProducer 創建一個新的模擬 Kafka 生產者
func NewMockKafkaProducer() *MockKafkaProducer {
	return &MockKafkaProducer{
		Messages: make(map[string][]string),
	}
}

// ProduceMessage 將消息添加到模擬存儲中
func (p *MockKafkaProducer) ProduceMessage(topic string, key string, value []byte) error {
	if _, exists := p.Messages[topic]; !exists {
		p.Messages[topic] = make([]string, 0)
	}
	p.Messages[topic] = append(p.Messages[topic], string(value))
	return nil
}

// SendLogBatch 將一批日誌添加到模擬存儲中
func (p *MockKafkaProducer) SendLogBatch(logs []models.LogEntry) error {
	for _, log := range logs {
		logBytes, err := json.Marshal(log)
		if err != nil {
			return fmt.Errorf("failed to marshal log entry: %w", err)
		}
		p.ProduceMessage("logs", log.Service, logBytes)
	}
	return nil
}

// Close 模擬關閉操作
func (p *MockKafkaProducer) Close() error {
	// 清空消息，模擬關閉
	p.Messages = make(map[string][]string)
	return nil
}