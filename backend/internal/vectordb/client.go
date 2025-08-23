package vectordb

import (
	"context"
)

// VectorDBClient interface for different vector database providers
type VectorDBClient interface {
	Upsert(ctx context.Context, collection string, vectors []Vector) error
	Search(ctx context.Context, collection string, query Vector, limit int) ([]SearchResult, error)
	Delete(ctx context.Context, collection string, ids []string) error
	Close() error
}

// Vector represents a vector with ID and metadata
type Vector struct {
	ID       string                 `json:"id"`
	Vector   []float64              `json:"vector"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SearchResult represents a search result with score
type SearchResult struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// MockVectorDBClient for testing
type MockVectorDBClient struct {
	vectors map[string]map[string]Vector
}

// NewVectorDBClient creates a new vector database client based on environment variables
func NewVectorDBClient() (VectorDBClient, error) {
	// For now, return mock client
	// TODO: Implement actual client selection based on environment variables
	return &MockVectorDBClient{
		vectors: make(map[string]map[string]Vector),
	}, nil
}

// NewMockVectorDBClient creates a new mock vector database client
func NewMockVectorDBClient() *MockVectorDBClient {
	return &MockVectorDBClient{
		vectors: make(map[string]map[string]Vector),
	}
}

// Upsert implementation for MockVectorDBClient
func (m *MockVectorDBClient) Upsert(ctx context.Context, collection string, vectors []Vector) error {
	if m.vectors[collection] == nil {
		m.vectors[collection] = make(map[string]Vector)
	}
	
	for _, vector := range vectors {
		m.vectors[collection][vector.ID] = vector
	}
	
	return nil
}

// Search implementation for MockVectorDBClient
func (m *MockVectorDBClient) Search(ctx context.Context, collection string, query Vector, limit int) ([]SearchResult, error) {
	if m.vectors[collection] == nil {
		return []SearchResult{}, nil
	}
	
	var results []SearchResult
	for _, vector := range m.vectors[collection] {
		// Simple cosine similarity mock
		score := calculateCosineSimilarity(query.Vector, vector.Vector)
		results = append(results, SearchResult{
			ID:       vector.ID,
			Score:    score,
			Metadata: vector.Metadata,
		})
	}
	
	// Sort by score (descending) and limit results
	if len(results) > limit {
		results = results[:limit]
	}
	
	return results, nil
}

// Delete implementation for MockVectorDBClient
func (m *MockVectorDBClient) Delete(ctx context.Context, collection string, ids []string) error {
	if m.vectors[collection] == nil {
		return nil
	}
	
	for _, id := range ids {
		delete(m.vectors[collection], id)
	}
	
	return nil
}

// Close implementation for MockVectorDBClient
func (m *MockVectorDBClient) Close() error {
	return nil
}

// calculateCosineSimilarity calculates cosine similarity between two vectors
func calculateCosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}
	
	var dotProduct, normA, normB float64
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	
	if normA == 0 || normB == 0 {
		return 0.0
	}
	
	return dotProduct / (normA * normB)
}
