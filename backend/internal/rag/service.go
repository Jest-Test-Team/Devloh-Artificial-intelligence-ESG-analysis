package rag

import (
	"context"
	"fmt"
	"time"

	"devlog-ai/internal/embedding"
	"devlog-ai/internal/vectordb"
)

// RAGService provides Retrieval-Augmented Generation functionality
type RAGService struct {
	embeddingClient embedding.EmbeddingClient
	vectorDBClient  vectordb.VectorDBClient
}

// NewRAGService creates a new RAG service instance
func NewRAGService(embeddingClient embedding.EmbeddingClient, vectorDBClient vectordb.VectorDBClient) *RAGService {
	return &RAGService{
		embeddingClient: embeddingClient,
		vectorDBClient:  vectorDBClient,
	}
}

// Query processes a natural language query and returns an answer
func (r *RAGService) Query(ctx context.Context, query string, collection string) (string, error) {
	startTime := time.Now()

	// Create embedding for the query
	queryEmbedding, err := r.embeddingClient.CreateEmbedding(query)
	if err != nil {
		return "", fmt.Errorf("failed to create query embedding: %w", err)
	}

	// Search for similar vectors
	searchVector := vectordb.Vector{
		Vector: queryEmbedding,
	}

	results, err := r.vectorDBClient.Search(ctx, collection, searchVector, 5)
	if err != nil {
		return "", fmt.Errorf("failed to search vector database: %w", err)
	}

	// Generate answer based on retrieved context
	answer := r.generateAnswer(query, results)

	// Log processing time
	processingTime := time.Since(startTime).Milliseconds()

	// TODO: Implement actual answer generation logic
	// For now, return a simple response
	return fmt.Sprintf("Based on the retrieved context, here's what I found: %s (Processed in %dms)", 
		answer, processingTime), nil
}

// IndexDocument indexes a document for later retrieval
func (r *RAGService) IndexDocument(ctx context.Context, collection string, documentID string, content string, metadata map[string]interface{}) error {
	// Create embedding for the document
	embedding, err := r.embeddingClient.CreateEmbedding(content)
	if err != nil {
		return fmt.Errorf("failed to create document embedding: %w", err)
	}

	// Store in vector database
	vector := vectordb.Vector{
		ID:       documentID,
		Vector:   embedding,
		Metadata: metadata,
	}

	if err := r.vectorDBClient.Upsert(ctx, collection, []vectordb.Vector{vector}); err != nil {
		return fmt.Errorf("failed to store document in vector database: %w", err)
	}

	return nil
}

// generateAnswer generates an answer based on retrieved context
func (r *RAGService) generateAnswer(query string, results []vectordb.SearchResult) string {
	if len(results) == 0 {
		return "I couldn't find any relevant information to answer your question."
	}

	// Simple answer generation - in a real implementation, this would use an LLM
	answer := fmt.Sprintf("Based on %d relevant documents, I can help with your query: '%s'. ", len(results), query)
	
	for i, result := range results {
		if i < 3 { // Limit to first 3 results
			answer += fmt.Sprintf("Document %d (relevance: %.2f) contains relevant information. ", i+1, result.Score)
		}
	}

	return answer
}

// Close closes the RAG service and its dependencies
func (r *RAGService) Close() error {
	if r.vectorDBClient != nil {
		return r.vectorDBClient.Close()
	}
	return nil
}
