package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CloudflareClient 定義了與 Cloudflare RAG 函數交互的接口
type CloudflareClient interface {
	Query(query string, options map[string]interface{}) (string, error)
}

// cloudflareClient 實現了 CloudflareClient 接口，用於實際的 Cloudflare API 調用
type cloudflareClient struct {
	apiEndpoint string
	apiKey      string
	httpClient  *http.Client
}

// NewCloudflareClient 創建一個新的 Cloudflare RAG 客戶端
func NewCloudflareClient(apiEndpoint, apiKey string) CloudflareClient {
	return &cloudflareClient{
		apiEndpoint: apiEndpoint,
		apiKey:      apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Query 向 Cloudflare RAG 函數發送查詢並返回響應
func (c *cloudflareClient) Query(query string, options map[string]interface{}) (string, error) {
	requestBody := map[string]interface{}{
		"query": query,
	}
	for k, v := range options {
		requestBody[k] = v
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiEndpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Cloudflare API returned non-OK status: %s", resp.Status)
	}

	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	// 假設 RAG 響應在 "response" 或 "result" 字段中
	if result, ok := responseBody["response"].(string); ok {
		return result, nil
	}
	if result, ok := responseBody["result"].(string); ok {
		return result, nil
	}

	return "", fmt.Errorf("unexpected response format from Cloudflare RAG API")
}