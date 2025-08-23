package vectordb

import (
	"errors"
	"os"
)

// VectorDBClient 是向量數據庫客戶端的接口
type VectorDBClient interface {
	// Store 將嵌入向量存儲到向量數據庫中
	Store(id string, vector []float32, metadata map[string]interface{}) error
	
	// Query 查詢與給定向量最相似的向量
	Query(vector []float32, limit int) ([]QueryResult, error)
	
	// Delete 從數據庫中刪除向量
	Delete(id string) error
}

// QueryResult 表示查詢結果
type QueryResult struct {
	ID       string                 // 向量的唯一標識符
	Score    float32                // 相似度分數
	Metadata map[string]interface{} // 與向量相關的元數據
}

// ClientType 表示向量數據庫客戶端的類型
type ClientType string

const (
	// PineconeClientType 是 Pinecone 向量數據庫客戶端的類型
	PineconeClientType ClientType = "pinecone"
	// WeaviateClientType 是 Weaviate 向量數據庫客戶端的類型
	WeaviateClientType ClientType = "weaviate"
	// MockClientType 是用於測試的模擬向量數據庫客戶端的類型
	MockClientType ClientType = "mock"
)

// NewVectorDBClient 創建一個新的向量數據庫客戶端
func NewVectorDBClient() (VectorDBClient, error) {
	// 從環境變量獲取向量數據庫類型
	dbType := os.Getenv("VECTOR_DB_TYPE")
	
	switch ClientType(dbType) {
	case PineconeClientType:
		apiKey := os.Getenv("PINECONE_API_KEY")
		environment := os.Getenv("PINECONE_ENVIRONMENT")
		indexName := os.Getenv("PINECONE_INDEX_NAME")
		
		if apiKey == "" || environment == "" || indexName == "" {
			return nil, errors.New("PINECONE_API_KEY, PINECONE_ENVIRONMENT, and PINECONE_INDEX_NAME are required for Pinecone")
		}
		
		return NewPineconeClient(apiKey, environment, indexName), nil
		
	case WeaviateClientType:
		endpoint := os.Getenv("WEAVIATE_ENDPOINT")
		apiKey := os.Getenv("WEAVIATE_API_KEY")
		className := os.Getenv("WEAVIATE_CLASS_NAME")
		
		if endpoint == "" || className == "" {
			return nil, errors.New("WEAVIATE_ENDPOINT and WEAVIATE_CLASS_NAME are required for Weaviate")
		}
		
		return NewWeaviateClient(endpoint, apiKey, className), nil
		
	default:
		// 默認使用模擬客戶端
		return NewMockVectorDBClient(), nil
	}
}

// NewMockVectorDBClient 創建一個新的模擬向量數據庫客戶端
func NewMockVectorDBClient() VectorDBClient {
	return &MockVectorDBClient{
		vectors: make(map[string]mockVectorEntry),
	}
}

// MockVectorDBClient 是一個用於測試的模擬向量數據庫客戶端
type MockVectorDBClient struct {
	vectors map[string]mockVectorEntry
}

type mockVectorEntry struct {
	vector   []float32
	metadata map[string]interface{}
}

// Store 將向量存儲到模擬數據庫中
func (c *MockVectorDBClient) Store(id string, vector []float32, metadata map[string]interface{}) error {
	c.vectors[id] = mockVectorEntry{
		vector:   vector,
		metadata: metadata,
	}
	return nil
}

// Query 查詢模擬數據庫中與給定向量最相似的向量
func (c *MockVectorDBClient) Query(vector []float32, limit int) ([]QueryResult, error) {
	if len(c.vectors) == 0 {
		return []QueryResult{}, nil
	}
	
	// 計算所有向量的相似度
	results := make([]QueryResult, 0, len(c.vectors))
	for id, entry := range c.vectors {
		score := cosineSimilarity(vector, entry.vector)
		results = append(results, QueryResult{
			ID:       id,
			Score:    score,
			Metadata: entry.metadata,
		})
	}
	
	// 按相似度排序
	sortQueryResults(results)
	
	// 限制結果數量
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}
	
	return results, nil
}

// Delete 從模擬數據庫中刪除向量
func (c *MockVectorDBClient) Delete(id string) error {
	delete(c.vectors, id)
	return nil
}

// cosineSimilarity 計算兩個向量之間的餘弦相似度
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}
	
	var dotProduct, normA, normB float32
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	
	if normA == 0 || normB == 0 {
		return 0
	}
	
	return dotProduct / (float32(normA) * float32(normB))
}

// sortQueryResults 按相似度降序排序查詢結果
func sortQueryResults(results []QueryResult) {
	// 使用簡單的冒泡排序（實際應用中可能使用更高效的排序算法）
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}