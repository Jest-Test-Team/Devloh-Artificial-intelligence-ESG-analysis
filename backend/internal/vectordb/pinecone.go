package vectordb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PineconeClient 實現了 VectorDBClient 接口，用於與 Pinecone 向量數據庫交互
type PineconeClient struct {
	apiKey      string
	environment string
	indexName   string
	httpClient  *http.Client
}

// NewPineconeClient 創建一個新的 Pinecone 客戶端
func NewPineconeClient(apiKey, environment, indexName string) *PineconeClient {
	return &PineconeClient{
		apiKey:      apiKey,
		environment: environment,
		indexName:   indexName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Store 將向量存儲到 Pinecone 中
func (c *PineconeClient) Store(id string, vector []float32, metadata map[string]interface{}) error {
	// 構建 Pinecone API 請求
	url := fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io/vectors/upsert", c.indexName, c.environment, c.environment)
	
	// 構建請求體
	requestBody := map[string]interface{}{
		"vectors": []map[string]interface{}{
			{
				"id":       id,
				"values":   vector,
				"metadata": metadata,
			},
		},
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	// 創建請求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// 設置標頭
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)
	
	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 檢查響應
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return fmt.Errorf("pinecone API error: %v", errorResponse)
	}
	
	return nil
}

// Query 查詢與給定向量最相似的向量
func (c *PineconeClient) Query(vector []float32, limit int) ([]QueryResult, error) {
	// 構建 Pinecone API 請求
	url := fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io/query", c.indexName, c.environment, c.environment)
	
	// 構建請求體
	requestBody := map[string]interface{}{
		"vector":    vector,
		"topK":      limit,
		"includeMetadata": true,
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	// 創建請求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// 設置標頭
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)
	
	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 檢查響應
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, fmt.Errorf("pinecone API error: %v", errorResponse)
	}
	
	// 解析響應
	var response struct {
		Matches []struct {
			ID       string                 `json:"id"`
			Score    float32                `json:"score"`
			Metadata map[string]interface{} `json:"metadata"`
		} `json:"matches"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// 轉換為查詢結果
	results := make([]QueryResult, len(response.Matches))
	for i, match := range response.Matches {
		results[i] = QueryResult{
			ID:       match.ID,
			Score:    match.Score,
			Metadata: match.Metadata,
		}
	}
	
	return results, nil
}

// Delete 從 Pinecone 中刪除向量
func (c *PineconeClient) Delete(id string) error {
	// 構建 Pinecone API 請求
	url := fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io/vectors/delete", c.indexName, c.environment, c.environment)
	
	// 構建請求體
	requestBody := map[string]interface{}{
		"ids": []string{id},
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	// 創建請求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// 設置標頭
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", c.apiKey)
	
	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 檢查響應
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return fmt.Errorf("pinecone API error: %v", errorResponse)
	}
	
	return nil
}