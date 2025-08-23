package vectordb

import (
	"testing"
)

func TestNewVectorDBClient(t *testing.T) {
	// 測試默認行為（使用環境變量設置的 mock 類型）
	client, err := NewVectorDBClient()
	if err != nil {
		t.Fatalf("Failed to create default Vector DB client: %v", err)
	}

	// 檢查是否是模擬客戶端
	if _, ok := client.(*MockVectorDBClient); !ok {
		t.Error("Expected MockVectorDBClient when VECTOR_DB_TYPE is set to 'mock'")
	}

	// 由於我們在 CI/CD 環境中運行測試，我們不能真正測試其他類型的客戶端
	// 在實際的集成測試中，你可能會使用測試容器或模擬
	t.Skip("Skipping other Vector DB client tests in CI environment")
}

// 測試模擬Vector DB客戶端
func TestMockVectorDBClient(t *testing.T) {
	client := NewMockVectorDBClient()

	// 測試Store
	id := "test-id"
	vector := []float32{0.1, 0.2, 0.3}
	metadata := map[string]interface{}{
		"text": "Test document",
		"source": "test",
	}

	err := client.Store(id, vector, metadata)
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
