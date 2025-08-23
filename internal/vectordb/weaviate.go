package vectordb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// weaviateClient is a client for the Weaviate vector database
type weaviateClient struct {
	config     Config
	httpClient *http.Client
}

// WeaviateClass represents a class in Weaviate
const WeaviateClass = "LogEntry"

// WeaviateObject represents an object in Weaviate
type WeaviateObject struct {
	Class      string                 `json:"class"`
	ID         string                 `json:"id,omitempty"`
	Vector     []float32              `json:"vector"`
	Properties map[string]interface{} `json:"properties"`
}

// WeaviateBatchResponse represents a batch response from Weaviate
type WeaviateBatchResponse struct {
	Objects []struct {
		Status string `json:"status"`
		ID     string `json:"id"`
	} `json:"objects"`
}

// WeaviateGraphQLRequest represents a GraphQL request to Weaviate
type WeaviateGraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// WeaviateGraphQLResponse represents a GraphQL response from Weaviate
type WeaviateGraphQLResponse struct {
	Data struct {
		Get struct {
			LogEntry []struct {
				ID         string                 `json:"id"`
				Properties map[string]interface{} `json:"properties"`
				Vector     []float32              `json:"vector"`
				Score      float32                `json:"_additional,omitempty"`
			} `json:"LogEntry"`
		} `json:"Get"`
	} `json:"data"`
}

// newWeaviateClient creates a new Weaviate client
func newWeaviateClient(config Config) (Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Weaviate API key is required")
	}
	if config.URL == "" {
		return nil, fmt.Errorf("Weaviate URL is required")
	}

	return &weaviateClient{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// Connect implements the Client interface
func (c *weaviateClient) Connect() error {
	// Check if the class exists, create it if it doesn't
	url := fmt.Sprintf("%s/v1/schema", c.config.URL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error connecting to Weaviate: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Weaviate: %s - %s", resp.Status, string(body))
	}

	// TODO: Check if the class exists and create it if it doesn't
	
	return nil
}

// Close implements the Client interface
func (c *weaviateClient) Close() error {
	// Nothing to close for HTTP client
	return nil
}

// UpsertVector implements the Client interface
func (c *weaviateClient) UpsertVector(id string, vector []float32, metadata map[string]interface{}) error {
	return c.UpsertVectors([]VectorEntry{
		{
			ID:       id,
			Vector:   vector,
			Metadata: metadata,
		},
	})
}

// UpsertVectors implements the Client interface
func (c *weaviateClient) UpsertVectors(vectors []VectorEntry) error {
	// Create the objects
	objects := make([]map[string]interface{}, len(vectors))
	for i, v := range vectors {
		objects[i] = map[string]interface{}{
			"class":      WeaviateClass,
			"id":         v.ID,
			"vector":     v.Vector,
			"properties": v.Metadata,
		}
	}
	
	// Create the request
	requestData := map[string]interface{}{
		"objects": objects,
	}
	
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}
	
	// Send the request
	url := fmt.Sprintf("%s/v1/batch/objects", c.config.URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error upserting vectors: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Weaviate: %s - %s", resp.Status, string(body))
	}
	
	return nil
}

// QueryNearest implements the Client interface
func (c *weaviateClient) QueryNearest(vector []float32, limit int) ([]QueryResult, error) {
	// Create the GraphQL query
	vectorJSON, err := json.Marshal(vector)
	if err != nil {
		return nil, fmt.Errorf("error marshaling vector: %w", err)
	}

	query := fmt.Sprintf(`
	{
		Get {
			LogEntry(
				nearVector: {
					vector: %s
					certainty: 0.7
				}
				limit: %d
			) {
				id
				properties
				_additional {
					certainty
				}
			}
		}
	}`, string(vectorJSON), limit)
	
	// Create the request
	requestData := WeaviateGraphQLRequest{
		Query: query,
	}
	
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}
	
	// Send the request
	url := fmt.Sprintf("%s/v1/graphql", c.config.URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying vectors: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from Weaviate: %s - %s", resp.Status, string(body))
	}
	
	// Parse the response
	var response WeaviateGraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	
	// Convert to QueryResult
	results := make([]QueryResult, len(response.Data.Get.LogEntry))
	for i, entry := range response.Data.Get.LogEntry {
		results[i] = QueryResult{
			ID:       entry.ID,
			Score:    entry.Score,
			Metadata: entry.Properties,
		}
	}
	
	return results, nil
}

// DeleteVector implements the Client interface
func (c *weaviateClient) DeleteVector(id string) error {
	url := fmt.Sprintf("%s/v1/objects/%s", c.config.URL, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error deleting vector: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Weaviate: %s - %s", resp.Status, string(body))
	}
	
	return nil
}

// DeleteVectors implements the Client interface
func (c *weaviateClient) DeleteVectors(ids []string) error {
	for _, id := range ids {
		if err := c.DeleteVector(id); err != nil {
			return err
		}
	}
	
	return nil
}