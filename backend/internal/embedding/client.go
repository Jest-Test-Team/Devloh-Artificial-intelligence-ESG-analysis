package embedding

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// Client 是一個文本嵌入生成客戶端的接口
type Client interface {
	// CreateEmbedding 生成文本的嵌入向量
	CreateEmbedding(text string) ([]float32, error)
}

// ClientType 表示嵌入客戶端的類型
type ClientType string

const (
	// OpenAIClientType 是 OpenAI 嵌入客戶端的類型
	OpenAIClientType ClientType = "openai"
	// GeminiClientType 是 Google Gemini 嵌入客戶端的類型
	GeminiClientType ClientType = "gemini"
	// MockClientType 是用於測試的模擬嵌入客戶端的類型
	MockClientType ClientType = "mock"
)

// NewEmbeddingClient 創建一個新的嵌入客戶端
func NewEmbeddingClient() (Client, error) {
	// 從環境變量獲取嵌入提供者
	provider := os.Getenv("EMBEDDING_PROVIDER")
	
	switch ClientType(provider) {
	case OpenAIClientType:
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY is required for OpenAI embedding provider")
		}
		return NewOpenAIClient(apiKey), nil
		
	case GeminiClientType:
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY is required for Gemini embedding provider")
		}
		return NewGeminiClient(apiKey), nil
		
	default:
		// 默認使用模擬客戶端
		return NewMockEmbeddingClient(), nil
	}
}

// OpenAIClient 是 OpenAI 嵌入 API 的客戶端
type OpenAIClient struct {
	apiKey string
	model  string
}

// NewOpenAIClient 創建一個新的 OpenAI 嵌入客戶端
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		model:  "text-embedding-3-small", // OpenAI 的默認嵌入模型
	}
}

// CreateEmbedding 使用 OpenAI API 生成文本的嵌入向量
func (c *OpenAIClient) CreateEmbedding(text string) ([]float32, error) {
	// 在實際實現中，這裡應該調用 OpenAI 的 API
	// 現在我們只是返回一個模擬的嵌入向量
	return createMockEmbedding(text, 1536), nil
}

// GeminiClient 是 Google Gemini 嵌入 API 的客戶端
type GeminiClient struct {
	apiKey string
	model  string
}

// NewGeminiClient 創建一個新的 Gemini 嵌入客戶端
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		model:  "embedding-001", // Gemini 的默認嵌入模型
	}
}

// CreateEmbedding 使用 Gemini API 生成文本的嵌入向量
func (c *GeminiClient) CreateEmbedding(text string) ([]float32, error) {
	// 在實際實現中，這裡應該調用 Gemini 的 API
	// 現在我們只是返回一個模擬的嵌入向量
	return createMockEmbedding(text, 768), nil
}

// MockEmbeddingClient 是一個用於測試的模擬嵌入客戶端
type MockEmbeddingClient struct{}

// NewMockEmbeddingClient 創建一個新的模擬嵌入客戶端
func NewMockEmbeddingClient() *MockEmbeddingClient {
	return &MockEmbeddingClient{}
}

// CreateEmbedding 生成一個確定性的模擬嵌入向量
func (c *MockEmbeddingClient) CreateEmbedding(text string) ([]float32, error) {
	// 使用固定維度的嵌入
	return createMockEmbedding(text, 384), nil
}

// createMockEmbedding 基於輸入文本生成確定性的模擬嵌入向量
func createMockEmbedding(text string, dimensions int) []float32 {
	// 使用文本的哈希生成確定性的嵌入向量
	hasher := sha256.New()
	hasher.Write([]byte(text))
	hash := hasher.Sum(nil)
	
	// 創建指定維度的嵌入向量
	embedding := make([]float32, dimensions)
	
	// 使用哈希填充嵌入向量
	for i := 0; i < dimensions; i++ {
		// 使用哈希的不同部分生成每個維度的值
		seedBytes := make([]byte, 4)
		for j := 0; j < 4; j++ {
			hashIndex := (i*4 + j) % len(hash)
			seedBytes[j] = hash[hashIndex]
		}
		
		// 將字節轉換為 uint32，然後歸一化到 [-1, 1] 範圍
		seed := binary.BigEndian.Uint32(seedBytes)
		normalized := (float32(seed) / float32(math.MaxUint32))*2 - 1
		
		embedding[i] = normalized
	}
	
	// 歸一化嵌入向量
	normalizeVector(embedding)
	
	return embedding
}

// normalizeVector 將向量歸一化為單位向量
func normalizeVector(vector []float32) {
	// 計算向量的模
	var sum float32
	for _, v := range vector {
		sum += v * v
	}
	magnitude := float32(math.Sqrt(float64(sum)))
	
	// 防止除以零
	if magnitude > 0 {
		for i := range vector {
			vector[i] /= magnitude
		}
	}
}