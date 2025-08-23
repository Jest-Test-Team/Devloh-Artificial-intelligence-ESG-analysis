package embedding

import (
	"testing"
)

func TestNewEmbeddingClient(t *testing.T) {
	// 測試默認行為（使用環境變量設置的 mock 類型）
	client, err := NewEmbeddingClient()
	if err != nil {
		t.Fatalf("Failed to create default embedding client: %v", err)
	}

	// 檢查是否是模擬客戶端
	if _, ok := client.(*MockEmbeddingClient); !ok {
		t.Error("Expected MockEmbeddingClient when EMBEDDING_PROVIDER is set to 'mock'")
	}

	// 由於我們在 CI/CD 環境中運行測試，我們不能真正測試其他類型的客戶端
	// 在實際的集成測試中，你可能會使用測試容器或模擬
	t.Skip("Skipping other embedding client tests in CI environment")
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
