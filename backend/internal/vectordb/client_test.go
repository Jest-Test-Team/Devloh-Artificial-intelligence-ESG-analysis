package vectordb

import (
	"os"
	"testing"
)

func TestNewVectorDBClient(t *testing.T) {
	// 保存原始環境變量
	origProvider := os.Getenv("VECTOR_DB_PROVIDER")
	origPineconeKey := os.Getenv("PINECONE_API_KEY")
	origPineconeEnv := os.Getenv("PINECONE_ENVIRONMENT")
	origWeaviateURL := os.Getenv("WEAVIATE_URL")
	origWeaviateKey := os.Getenv("WEAVIATE_API_KEY")

	// 測試完成後恢復環境變量
	defer func() {
		os.Setenv("VECTOR_DB_PROVIDER", origProvider)
		os.Setenv("PINECONE_API_KEY", origPineconeKey)
		os.Setenv("PINECONE_ENVIRONMENT", origPineconeEnv)
		os.Setenv("WEAVIATE_URL", origWeaviateURL)
		os.Setenv("WEAVIATE_API_KEY", origWeaviateKey)
	}()

	// 測試默認行為（沒有設置環境變量）
	os.Unsetenv("VECTOR_DB_PROVIDER")
	os.Unsetenv("PINECONE_API_KEY")
	os.Unsetenv("PINECONE_ENVIRONMENT")
	os.Unsetenv("WEAVIATE_URL")
	os.Unsetenv("WEAVIATE_API_KEY")

	client, err := NewVectorDBClient()
	if err != nil {
		t.Fatalf("Failed to create default Vector DB client: %v", err)
	}

	// 檢查是否是模擬客戶端
	if _, ok := client.(*MockVectorDBClient); !ok {
		t.Error("Expected MockVectorDBClient when environment variables are not set")
	}

	// 測試Pinecone客戶端
	os.Setenv("VECTOR_DB_PROVIDER", "pinecone")
	os.Setenv("PINECONE_API_KEY", "test-api-key")
	os.Setenv("PINECONE_ENVIRONMENT", "test-environment")

	client, err = NewVectorDBClient()
	if err != nil {
		t.Fatalf("Failed to create Pinecone client: %v", err)
	}

	// 檢查是否是Pinecone客戶端
	if _, ok := client.(*PineconeClient); !ok {
		t.Error("Expected PineconeClient when VECTOR_DB_PROVIDER is set to 'pinecone'")
	}

	// 測試Weaviate客戶端
	os.Setenv("VECTOR_DB_PROVIDER", "weaviate")
	os.Setenv("WEAVIATE_URL", "http://localhost:8080")
	os.Setenv("WEAVIATE_API_KEY", "test-api-key")

	client, err = NewVectorDBClient()
	if err != nil {
		t.Fatalf("Failed to create Weaviate client: %v", err)
	}

	// 檢查是否是Weaviate客戶端
	if _, ok := client.(*WeaviateClient); !ok {
		t.Error("Expected WeaviateClient when VECTOR_DB_PROVIDER is set to 'weaviate'")
	}

	// 測試無效的提供者
	os.Setenv("VECTOR_DB_PROVIDER", "invalid")

	_, err = NewVectorDBClient()
	if err == nil {
		t.Error("Expected error when VECTOR_DB_PROVIDER is invalid")
	}
}

// 測試模擬Vector DB客戶端
func TestMockVectorDBClient(t *testing.T) {
	client := NewMockVectorDBClient()

	// 測試Upsert
	id := "test-id"
	vector := []float32{0.1, 0.2, 0.3}
	metadata := map[string]interface{}{
		"text": "Test document",
		"source": "test",
	}

	err := client.Upsert(id, vector, metadata)
	if err != nil {
		t.Fatalf("Failed to upsert document: %v", err)
	}

	// 測試Query
	results, err := client.Query(vector, 5)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0].ID != id {
		t.Errorf("Expected ID '%s', got '%s'", id, results[0].ID)
	}

	// 測試Delete
	err = client.Delete(id)
	if err != nil {
		t.Fatalf("Failed to delete document: %v", err)
	}

	// 刪除後查詢應該返回空結果
	results, err = client.Query(vector, 5)
	if err != nil {
		t.Fatalf("Failed to query after delete: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results after delete, got %d", len(results))
	}
}
