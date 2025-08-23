package vectordb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// PineconeClient provides access to Pinecone vector database
type PineconeClient struct {
	apiKey      string
	environment string
	indexName   string
	baseURL     string
	client      *http.Client
}

// PineconeUpsertRequest represents an upsert request to Pinecone
type PineconeUpsertRequest struct {
	Vectors []PineconeVector `json:"vectors"`
}

// PineconeVector represents a vector in Pinecone
type PineconeVector struct {
	ID       string                 `json:"id"`
	Values   []float64              `json:"values"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// PineconeQueryRequest represents a query request to Pinecone
type PineconeQueryRequest struct {
	Vector          []float64              `json:"vector"`
	TopK            int                    `json:"topK"`
	IncludeMetadata bool                   `json:"includeMetadata"`
	Filter          map[string]interface{} `json:"filter,omitempty"`
}

// PineconeQueryResponse represents a query response from Pinecone
type PineconeQueryResponse struct {
	Matches []PineconeMatch `json:"matches"`
}

// PineconeMatch represents a match from Pinecone
type PineconeMatch struct {
	ID       string                 `json:"id"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

// NewPineconeClient creates a new Pinecone client
func NewPineconeClient(apiKey, environment, indexName string) *PineconeClient {
	return &PineconeClient{
		apiKey:      apiKey,
		environment: environment,
		indexName:   indexName,
		baseURL:     fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io", indexName, environment, environment),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Upsert upserts vectors to Pinecone
func (p *PineconeClient) Upsert(ctx context.Context, collection string, vectors []Vector) error {
	pineconeVectors := make([]PineconeVector, len(vectors))
	for i, v := range vectors {
		pineconeVectors[i] = PineconeVector{
			ID:       v.ID,
			Values:   v.Vector,
			Metadata: v.Metadata,
		}
	}

	request := PineconeUpsertRequest{
		Vectors: pineconeVectors,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/vectors/upsert", p.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Api-Key", p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Search searches for similar vectors in Pinecone
func (p *PineconeClient) Search(ctx context.Context, collection string, query Vector, limit int) ([]SearchResult, error) {
	request := PineconeQueryRequest{
		Vector:          query.Vector,
		TopK:            limit,
		IncludeMetadata: true,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/query", p.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Api-Key", p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response PineconeQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	results := make([]SearchResult, len(response.Matches))
	for i, match := range response.Matches {
		results[i] = SearchResult{
			ID:       match.ID,
			Score:    match.Score,
			Metadata: match.Metadata,
		}
	}

	return results, nil
}

// Delete deletes vectors from Pinecone
func (p *PineconeClient) Delete(ctx context.Context, collection string, ids []string) error {
	// Pinecone delete endpoint implementation
	// TODO: Implement actual delete logic
	return nil
}

// Close closes the Pinecone client
func (p *PineconeClient) Close() error {
	// No resources to close for HTTP client
	return nil
}
