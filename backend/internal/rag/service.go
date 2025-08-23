package rag

import (
	"devlog-ai/internal/embedding"
	"devlog-ai/internal/vectordb"
	"fmt"
)

// Service 是 RAG（檢索增強生成）服務的接口
type Service interface {
	// Query 處理用戶查詢並返回增強的回應
	Query(query string, options map[string]interface{}) (string, error)
}

// service 實現了 RAG 服務接口
type service struct {
	embeddingClient embedding.Client
	vectorDBClient  vectordb.VectorDBClient
	cloudflareClient CloudflareClient
}

// NewService 創建一個新的 RAG 服務
func NewService(embeddingClient embedding.Client, vectorDBClient vectordb.VectorDBClient, cloudflareClient CloudflareClient) Service {
	return &service{
		embeddingClient: embeddingClient,
		vectorDBClient:  vectorDBClient,
		cloudflareClient: cloudflareClient,
	}
}

// Query 處理用戶查詢，檢索相關上下文，並使用 Cloudflare RAG 生成回應
func (s *service) Query(query string, options map[string]interface{}) (string, error) {
	// 檢查查詢是否為空
	if query == "" {
		return "", fmt.Errorf("query cannot be empty")
	}
	
	// 為查詢生成嵌入向量
	queryEmbedding, err := s.embeddingClient.CreateEmbedding(query)
	if err != nil {
		return "", fmt.Errorf("failed to create embedding for query: %w", err)
	}
	
	// 從向量數據庫中檢索相關上下文
	limit := 5 // 默認檢索 5 個相似文檔
	if limitOption, ok := options["limit"].(float64); ok {
		limit = int(limitOption)
	}
	
	results, err := s.vectorDBClient.Query(queryEmbedding, limit)
	if err != nil {
		return "", fmt.Errorf("failed to query vector database: %w", err)
	}
	
	// 準備上下文
	context := prepareContext(results)
	
	// 將上下文添加到選項中
	if options == nil {
		options = make(map[string]interface{})
	}
	options["context"] = context
	
	// 使用 Cloudflare RAG 處理查詢
	response, err := s.cloudflareClient.Query(query, options)
	if err != nil {
		return "", fmt.Errorf("failed to process query with Cloudflare RAG: %w", err)
	}
	
	return response, nil
}

// prepareContext 從查詢結果中提取上下文
func prepareContext(results []vectordb.QueryResult) string {
	context := ""
	
	for _, result := range results {
		// 從元數據中提取內容
		if content, ok := result.Metadata["content"].(string); ok {
			context += content + "\n\n"
		}
	}
	
	return context
}