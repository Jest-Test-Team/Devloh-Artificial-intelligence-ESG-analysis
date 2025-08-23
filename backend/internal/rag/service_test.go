package rag

import (
	"testing"

	"devlog-ai/internal/embedding"
	"devlog-ai/internal/vectordb"
)

// 模擬CloudflareRAG客戶端
type mockCloudflareRAGClient struct{}

func (m *mockCloudflareRAGClient) Query(query string, options map[string]interface{}) (string, error) {
	return "這是來自Cloudflare RAG的模擬回答", nil
}

// 測試RAG服務
func TestRAGService_Query(t *testing.T) {
	// 創建模擬客戶端
	mockCloudflare := &mockCloudflareRAGClient{}
	mockEmbedding := embedding.NewMockEmbeddingClient()
	mockVectorDB := vectordb.NewMockVectorDBClient()

	// 預先插入一些向量數據
	vector := []float32{0.1, 0.2, 0.3}
	metadata1 := map[string]interface{}{
		"content":   "這是第一條測試日誌",
		"timestamp": "2023-07-15T12:00:00Z",
		"service":   "test-service",
		"level":     "INFO",
	}
	metadata2 := map[string]interface{}{
		"content":   "這是第二條測試日誌",
		"timestamp": "2023-07-15T12:01:00Z",
		"service":   "test-service",
		"level":     "WARN",
	}

	mockVectorDB.Store("log-1", vector, metadata1)
	mockVectorDB.Store("log-2", vector, metadata2)

	// 創建RAG服務
	service := NewService(mockEmbedding, mockVectorDB, mockCloudflare)

	// 測試查詢
	query := "查詢系統錯誤"
	options := map[string]interface{}{
		"max_results":     5,
		"time_range_hours": 24,
	}

	response, err := service.Query(query, options)
	if err != nil {
		t.Fatalf("Failed to query RAG service: %v", err)
	}

	// 檢查響應
	if response != "這是來自Cloudflare RAG的模擬回答" {
		t.Errorf("Unexpected answer: %s", response)
	}
}

// 測試RAG服務的邊緣情況
func TestRAGService_EdgeCases(t *testing.T) {
	// 創建模擬客戶端
	mockCloudflare := &mockCloudflareRAGClient{}
	mockEmbedding := embedding.NewMockEmbeddingClient()
	mockVectorDB := vectordb.NewMockVectorDBClient()

	// 創建RAG服務
	service := NewService(mockEmbedding, mockVectorDB, mockCloudflare)

	// 測試空查詢
	t.Run("Empty Query", func(t *testing.T) {
		query := ""
		options := map[string]interface{}{
			"max_results": 5,
		}

		_, err := service.Query(query, options)
		if err == nil {
			t.Error("Expected error for empty query, got nil")
		}
	})

	// 測試選項驗證
	t.Run("Options Validation", func(t *testing.T) {
		query := "測試查詢"
		
		// 無效的max_results
		options := map[string]interface{}{
			"max_results": -1,
		}

		_, err := service.Query(query, options)
		if err != nil {
			// 我們現在允許負數的max_results，它會被忽略
			t.Errorf("Unexpected error for invalid max_results: %v", err)
		}
	})
}