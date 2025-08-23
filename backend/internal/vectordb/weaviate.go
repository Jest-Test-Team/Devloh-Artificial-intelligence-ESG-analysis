package vectordb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WeaviateClient provides access to Weaviate vector database
type WeaviateClient struct {
	baseURL string
	client  *http.Client
}

// WeaviateObject represents an object in Weaviate
type WeaviateObject struct {
	Class      string                 `json:"class"`
	ID         string                 `json:"id,omitempty"`
	Properties map[string]interface{} `json:"properties"`
	Vector     []float64              `json:"vector,omitempty"`
}

// WeaviateQuery represents a query to Weaviate
type WeaviateQuery struct {
	Class      string                 `json:"class"`
	Properties []string               `json:"properties"`
	NearVector map[string]interface{} `json:"nearVector"`
	Limit      int                    `json:"limit"`
}

// WeaviateResponse represents a response from Weaviate
type WeaviateResponse struct {
	Data struct {
		Get map[string][]WeaviateObject `json:"Get"`
	} `json:"data"`
}

// NewWeaviateClient creates a new Weaviate client
func NewWeaviateClient(baseURL string) *WeaviateClient {
	return &WeaviateClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Upsert upserts objects to Weaviate
func (w *WeaviateClient) Upsert(ctx context.Context, collection string, vectors []Vector) error {
	for _, v := range vectors {
		object := WeaviateObject{
			Class: collection,
			ID:    v.ID,
			Properties: map[string]interface{}{
				"id":       v.ID,
				"metadata": v.Metadata,
			},
			Vector: v.Vector,
		}

		requestBody, err := json.Marshal(object)
		if err != nil {
			return fmt.Errorf("failed to marshal object: %w", err)
		}

		url := fmt.Sprintf("%s/v1/objects", w.baseURL)
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := w.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
		}
	}

	return nil
}

// Search searches for similar objects in Weaviate
func (w *WeaviateClient) Search(ctx context.Context, collection string, query Vector, limit int) ([]SearchResult, error) {
	queryObj := WeaviateQuery{
		Class: collection,
		Properties: []string{
			"id",
			"metadata",
		},
		NearVector: map[string]interface{}{
			"vector": query.Vector,
		},
		Limit: limit,
	}

	requestBody, err := json.Marshal(map[string]interface{}{
		"query": queryObj,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	url := fmt.Sprintf("%s/v1/graphql", w.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response WeaviateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	objects := response.Data.Get[collection]
	results := make([]SearchResult, len(objects))
	for i, obj := range objects {
		results[i] = SearchResult{
			ID:       obj.ID,
			Score:    0.0, // Weaviate doesn't provide scores in this format
			Metadata: obj.Properties,
		}
	}

	return results, nil
}

// Delete deletes objects from Weaviate
func (w *WeaviateClient) Delete(ctx context.Context, collection string, ids []string) error {
	for _, id := range ids {
		url := fmt.Sprintf("%s/v1/objects/%s/%s", w.baseURL, collection, id)
		req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := w.client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			return fmt.Errorf("delete request failed with status: %d", resp.StatusCode)
		}
	}

	return nil
}

// Close closes the Weaviate client
func (w *WeaviateClient) Close() error {
	// No resources to close for HTTP client
	return nil
}
