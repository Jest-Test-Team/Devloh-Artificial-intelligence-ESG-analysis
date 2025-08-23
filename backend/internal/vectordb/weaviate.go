package vectordb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeaviateClient 實現了 VectorDBClient 接口，用於與 Weaviate 向量數據庫交互
type WeaviateClient struct {
	endpoint   string
	apiKey     string
	className  string
	httpClient *http.Client
}

// NewWeaviateClient 創建一個新的 Weaviate 客戶端
func NewWeaviateClient(endpoint, apiKey, className string) *WeaviateClient {
	return &WeaviateClient{
		endpoint:  endpoint,
		apiKey:    apiKey,
		className: className,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Store 將向量存儲到 Weaviate 中
func (c *WeaviateClient) Store(id string, vector []float32, metadata map[string]interface{}) error {
	// 構建 Weaviate API 請求
	url := fmt.Sprintf("%s/v1/objects", c.endpoint)
	
	// 構建請求體
	requestBody := map[string]interface{}{
		"class": c.className,
		"id":    id,
		"properties": metadata,
		"vector": vector,
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
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	
	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 檢查響應
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return fmt.Errorf("weaviate API error: %v", errorResponse)
	}
	
	return nil
}

// Query 查詢與給定向量最相似的向量
func (c *WeaviateClient) Query(vector []float32, limit int) ([]QueryResult, error) {
	// 構建 Weaviate API 請求
	url := fmt.Sprintf("%s/v1/graphql", c.endpoint)
	
	// 構建 GraphQL 查詢
	graphqlQuery := fmt.Sprintf(`
	{
		Get {
			%s(
				nearVector: {
					vector: %v,
					certainty: 0.7
				},
				limit: %d
			) {
				_additional {
					id
					certainty
				}
				%s
			}
		}
	}
	`, c.className, vectorToString(vector), limit, buildPropertiesQuery(getPropertyNames(map[string]interface{}{})))
	
	// 構建請求體
	requestBody := map[string]interface{}{
		"query": graphqlQuery,
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
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	
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
		return nil, fmt.Errorf("weaviate API error: %v", errorResponse)
	}
	
	// 解析響應
	var response struct {
		Data struct {
			Get map[string][]struct {
				Additional struct {
					ID        string  `json:"id"`
					Certainty float32 `json:"certainty"`
				} `json:"_additional"`
				Properties map[string]interface{} `json:"-"`
			} `json:"Get"`
		} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	// 轉換為查詢結果
	results := []QueryResult{}
	for _, item := range response.Data.Get[c.className] {
		results = append(results, QueryResult{
			ID:       item.Additional.ID,
			Score:    item.Additional.Certainty,
			Metadata: item.Properties,
		})
	}
	
	return results, nil
}

// Delete 從 Weaviate 中刪除向量
func (c *WeaviateClient) Delete(id string) error {
	// 構建 Weaviate API 請求
	url := fmt.Sprintf("%s/v1/objects/%s", c.endpoint, id)
	
	// 創建請求
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// 設置標頭
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}
	
	// 發送請求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 檢查響應
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return fmt.Errorf("weaviate API error: %v", errorResponse)
	}
	
	return nil
}

// vectorToString 將向量轉換為字符串表示
func vectorToString(vector []float32) string {
	result := "["
	for i, v := range vector {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%f", v)
	}
	result += "]"
	return result
}

// getPropertyNames 獲取元數據的屬性名稱
func getPropertyNames(metadata map[string]interface{}) []string {
	names := make([]string, 0, len(metadata))
	for name := range metadata {
		names = append(names, name)
	}
	return names
}

// buildPropertiesQuery 構建 GraphQL 屬性查詢
func buildPropertiesQuery(properties []string) string {
	if len(properties) == 0 {
		return "content timestamp source"
	}
	
	result := ""
	for _, prop := range properties {
		result += prop + " "
	}
	return result
}