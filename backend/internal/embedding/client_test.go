package embedding

import (
	"os"
	"testing"
)

func TestNewEmbeddingClient(t *testing.T) {
	// 保存原始環境變量
	origProvider := os.Getenv("EMBEDDING_PROVIDER")
	origOpenAIKey := os.Getenv("OPENAI_API_KEY")
	origGeminiKey := os.Getenv("GEMINI_API_KEY")

	// 測試完成後恢復環境變量
	defer func() {
		os.Setenv("EMBEDDING_PROVIDER", origProvider)
		os.Setenv("OPENAI_API_KEY", origOpenAIKey)
		os.Setenv("GEMINI_API_KEY", origGeminiKey)
	}()

	// 測試默認行為（沒有設置環境變量）
	os.Unsetenv("EMBEDDING_PROVIDER")
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")

	client, err := NewEmbeddingClient()
	if err != nil {
		t.Fatalf("Failed to create default embedding client: %v", err)
	}

	// 檢查是否是模擬客戶端
	if _, ok := client.(*MockEmbeddingClient); !ok {
		t.Error("Expected MockEmbeddingClient when environment variables are not set")
	}

	// 測試OpenAI客戶端
	os.Setenv("EMBEDDING_PROVIDER", "openai")
	os.Setenv("OPENAI_API_KEY", "test-api-key")

	client, err = NewEmbeddingClient()
	if err != nil {
		t.Fatalf("Failed to create OpenAI client: %v", err)
	}

	// 檢查是否是OpenAI客戶端
	if _, ok := client.(*OpenAIClient); !ok {
		t.Error("Expected OpenAIClient when EMBEDDING_PROVIDER is set to 'openai'")
	}

	// 測試Gemini客戶端
	os.Setenv("EMBEDDING_PROVIDER", "gemini")
	os.Setenv("GEMINI_API_KEY", "test-api-key")

	client, err = NewEmbeddingClient()
	if err != nil {
		t.Fatalf("Failed to create Gemini client: %v", err)
	}

	// 檢查是否是Gemini客戶端
	if _, ok := client.(*GeminiClient); !ok {
		t.Error("Expected GeminiClient when EMBEDDING_PROVIDER is set to 'gemini'")
	}

	// 測試無效的提供者
	os.Setenv("EMBEDDING_PROVIDER", "invalid")

	_, err = NewEmbeddingClient()
	if err == nil {
		t.Error("Expected error when EMBEDDING_PROVIDER is invalid")
	}
}

// 測試模擬嵌入客戶端
func TestMockEmbeddingClient(t *testing.T) {
	client := NewMockEmbeddingClient()

	// 測試創建嵌入
	text := "This is a test document"
	embedding, err := client.CreateEmbedding(text)
	if err != nil {
		t.Fatalf("Failed to create embedding: %v", err)
	}

	// 檢查嵌入維度
	if len(embedding) != 384 {
		t.Errorf("Expected embedding dimension 384, got %d", len(embedding))
	}

	// 檢查嵌入值是否在合理範圍內
	for i, value := range embedding {
		if value < -1.0 || value > 1.0 {
			t.Errorf("Embedding value at index %d is out of range: %f", i, value)
		}
	}

	// 測試相同文本產生相同的嵌入
	embedding2, err := client.CreateEmbedding(text)
	if err != nil {
		t.Fatalf("Failed to create second embedding: %v", err)
	}

	for i := range embedding {
		if embedding[i] != embedding2[i] {
			t.Errorf("Expected identical embeddings for same text, but they differ at index %d", i)
			break
		}
	}

	// 測試不同文本產生不同的嵌入
	differentText := "This is a completely different document"
	differentEmbedding, err := client.CreateEmbedding(differentText)
	if err != nil {
		t.Fatalf("Failed to create different embedding: %v", err)
	}

	// 檢查至少有一個值不同
	allSame := true
	for i := range embedding {
		if embedding[i] != differentEmbedding[i] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Error("Expected different embeddings for different texts, but they are identical")
	}
}
