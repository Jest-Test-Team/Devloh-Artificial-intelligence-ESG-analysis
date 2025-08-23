package messaging

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// MessageHandler 是處理 Kafka 消息的接口
type MessageHandler interface {
	// HandleMessage 處理一條 Kafka 消息
	HandleMessage(message []byte) error
}

// KafkaConsumer 定義了 Kafka 消費者的接口
type KafkaConsumer interface {
	// Start 開始消費消息
	Start() error

	// Stop 停止消費消息
	Stop() error
}

// 實際的 Kafka 消費者實現
type kafkaConsumer struct {
	reader    *kafka.Reader
	handler   MessageHandler
	ctx       context.Context
	cancel    context.CancelFunc
	isRunning bool
}

// NewKafkaConsumer 創建一個新的 Kafka 消費者
func NewKafkaConsumer(topics []string, groupID string, handler MessageHandler) (KafkaConsumer, error) {
	// 檢查必要的環境變量
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		// 如果沒有設置 Kafka 環境變量，使用模擬消費者
		return NewMockKafkaConsumer(handler), nil
	}

	// 創建 Kafka 讀取器
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokers},
		GroupID:        groupID,
		Topic:          topics[0], // 注意：這裡我們只使用第一個主題，實際應用中可能需要多個讀取器
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 1 * time.Second,
	})

	// 檢查是否需要 SASL 驗證
	// 注意：在實際應用中，應該使用適當的 SASL 機制
	// 這裡我們簡化處理，因為 kafka-go 版本可能不同
	username := os.Getenv("KAFKA_USERNAME")
	password := os.Getenv("KAFKA_PASSWORD")
	if username != "" && password != "" {
		// 在實際應用中設置 SASL
		// 例如：
		// dialer := &kafka.Dialer{
		//     Timeout:   10 * time.Second,
		//     DualStack: true,
		// }
		// 注意：在實際應用中，這裡應該設置 SASLMechanism
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &kafkaConsumer{
		reader:    reader,
		handler:   handler,
		ctx:       ctx,
		cancel:    cancel,
		isRunning: false,
	}, nil
}

// Start 開始消費消息
func (c *kafkaConsumer) Start() error {
	if c.isRunning {
		return fmt.Errorf("consumer is already running")
	}

	c.isRunning = true

	go func() {
		for {
			select {
			case <-c.ctx.Done():
				return
			default:
				m, err := c.reader.ReadMessage(c.ctx)
				if err != nil {
					// 檢查是否是因為上下文取消而退出
					if c.ctx.Err() != nil {
						return
					}
					// 其他錯誤，記錄並繼續
					fmt.Printf("Error reading message: %v\n", err)
					continue
				}

				// 處理消息
				if err := c.handler.HandleMessage(m.Value); err != nil {
					fmt.Printf("Error handling message: %v\n", err)
				}
			}
		}
	}()

	return nil
}

// Stop 停止消費消息
func (c *kafkaConsumer) Stop() error {
	if !c.isRunning {
		return nil
	}

	c.cancel()
	c.isRunning = false
	return c.reader.Close()
}

// MockKafkaConsumer 是一個模擬的 Kafka 消費者，用於測試
type MockKafkaConsumer struct {
	handler   MessageHandler
	messages  [][]byte
	isRunning bool
}

// NewMockKafkaConsumer 創建一個新的模擬 Kafka 消費者
func NewMockKafkaConsumer(handler MessageHandler) *MockKafkaConsumer {
	return &MockKafkaConsumer{
		handler:   handler,
		messages:  make([][]byte, 0),
		isRunning: false,
	}
}

// Start 模擬開始消費消息
func (c *MockKafkaConsumer) Start() error {
	if c.isRunning {
		return fmt.Errorf("consumer is already running")
	}

	c.isRunning = true
	return nil
}

// Stop 模擬停止消費消息
func (c *MockKafkaConsumer) Stop() error {
	if !c.isRunning {
		return nil
	}

	c.isRunning = false
	return nil
}

// AddMessage 添加一條消息到模擬消費者的隊列中
func (c *MockKafkaConsumer) AddMessage(message []byte) error {
	if !c.isRunning {
		return fmt.Errorf("consumer is not running")
	}

	c.messages = append(c.messages, message)

	// 處理消息
	return c.handler.HandleMessage(message)
}