package messaging

import (
	"encoding/json"
	"os"
	"testing"

	"devlog-ai/internal/api/models"
)

// 測試用的模擬Kafka生產者
type testKafkaProducer struct {
	messages map[string][]string // topic -> messages
}

func newTestKafkaProducer() *testKafkaProducer {
	return &testKafkaProducer{
		messages: make(map[string][]string),
	}
}

func (p *testKafkaProducer) ProduceMessage(topic string, key string, value []byte) error {
	if _, exists := p.messages[topic]; !exists {
		p.messages[topic] = make([]string, 0)
	}
	p.messages[topic] = append(p.messages[topic], string(value))
	return nil
}

func (p *testKafkaProducer) SendLogBatch(logs []models.LogEntry) error {
	for _, log := range logs {
		logBytes, _ := json.Marshal(log)
		p.ProduceMessage("logs", "", logBytes)
	}
	return nil
}

func (p *testKafkaProducer) Close() error {
	return nil
}

func TestNewKafkaProducer(t *testing.T) {
	// 保存原始環境變量
	origBroker := os.Getenv("KAFKA_BROKER")
	origUsername := os.Getenv("KAFKA_USERNAME")
	origPassword := os.Getenv("KAFKA_PASSWORD")

	// 測試完成後恢復環境變量
	defer func() {
		os.Setenv("KAFKA_BROKER", origBroker)
		os.Setenv("KAFKA_USERNAME", origUsername)
		os.Setenv("KAFKA_PASSWORD", origPassword)
	}()

	// 測試沒有設置環境變量時（應該返回模擬生產者）
	os.Unsetenv("KAFKA_BROKERS")
	os.Unsetenv("KAFKA_USERNAME")
	os.Unsetenv("KAFKA_PASSWORD")

	producer, err := NewKafkaProducer()
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}

	// 檢查是否是模擬生產者
	if _, ok := producer.(*MockKafkaProducer); !ok {
		t.Error("Expected MockKafkaProducer when environment variables are not set")
	}

	// 測試設置了環境變量但連接不可用時（這裡我們無法真正測試連接，但可以檢查代碼路徑）
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("KAFKA_USERNAME", "test-user")
	os.Setenv("KAFKA_PASSWORD", "test-password")

	// 在這裡我們不期望失敗，因為我們的實現會在連接失敗時仍然返回一個有效的生產者
	producer, err = NewKafkaProducer()
	if err != nil {
		t.Errorf("Unexpected error when creating Kafka producer: %v", err)
	}
}

func TestMockKafkaProducer(t *testing.T) {
	producer := NewMockKafkaProducer()

	// 測試發送單條消息
	err := producer.ProduceMessage("test-topic", "test-key", []byte("test-message"))
	if err != nil {
		t.Fatalf("Failed to produce message: %v", err)
	}

	// 測試發送日誌批次
	logs := []models.LogEntry{
		{
			Service:   "test-service-1",
			Level:     "INFO",
			Message:   "Test message 1",
			Timestamp: "2023-07-15T12:00:00Z",
		},
		{
			Service:   "test-service-2",
			Level:     "WARN",
			Message:   "Test message 2",
			Timestamp: "2023-07-15T12:01:00Z",
		},
	}

	err = producer.SendLogBatch(logs)
	if err != nil {
		t.Fatalf("Failed to send log batch: %v", err)
	}

	// 測試關閉
	err = producer.Close()
	if err != nil {
		t.Fatalf("Failed to close producer: %v", err)
	}
}
