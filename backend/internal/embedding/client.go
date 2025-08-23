package embedding

import (
	"fmt"
	"os"
)

// EmbeddingClient interface for different embedding providers
type EmbeddingClient interface {
	CreateEmbedding(text string) ([]float64, error)
}

// MockEmbeddingClient for testing
type MockEmbeddingClient struct{}

// OpenAIClient for OpenAI embeddings
type OpenAIClient struct {
	apiKey string
}

// GeminiClient for Google Gemini embeddings
type GeminiClient struct {
	apiKey string
}

// NewEmbeddingClient creates a new embedding client based on environment variables
func NewEmbeddingClient() (EmbeddingClient, error) {
	provider := os.Getenv("EMBEDDING_PROVIDER")
	
	switch provider {
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		return &OpenAIClient{apiKey: apiKey}, nil
	case "gemini":
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY not set")
		}
		return &GeminiClient{apiKey: apiKey}, nil
	default:
		// Return mock client for testing/development
		return &MockEmbeddingClient{}, nil
	}
}

// NewMockEmbeddingClient creates a new mock embedding client
func NewMockEmbeddingClient() *MockEmbeddingClient {
	return &MockEmbeddingClient{}
}

// CreateEmbedding implementation for MockEmbeddingClient
func (m *MockEmbeddingClient) CreateEmbedding(text string) ([]float64, error) {
	// Return a mock 384-dimensional embedding
	embedding := make([]float64, 384)
	for i := range embedding {
		embedding[i] = 0.1 * float64(i%10) - 0.5 // Simple mock values
	}
	return embedding, nil
}

// CreateEmbedding implementation for OpenAIClient
func (o *OpenAIClient) CreateEmbedding(text string) ([]float64, error) {
	// TODO: Implement actual OpenAI API call
	// For now, return mock data
	embedding := make([]float64, 1536) // OpenAI embeddings are typically 1536-dimensional
	for i := range embedding {
		embedding[i] = 0.1 * float64(i%10) - 0.5
	}
	return embedding, nil
}

// CreateEmbedding implementation for GeminiClient
func (g *GeminiClient) CreateEmbedding(text string) ([]float64, error) {
	// TODO: Implement actual Gemini API call
	// For now, return mock data
	embedding := make([]float64, 768) // Gemini embeddings are typically 768-dimensional
	for i := range embedding {
		embedding[i] = 0.1 * float64(i%10) - 0.5
	}
	return embedding, nil
}
