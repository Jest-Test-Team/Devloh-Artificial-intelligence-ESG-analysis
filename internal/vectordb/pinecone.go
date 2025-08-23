package vectordb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// pineconeClient is a client for the Pinecone vector database
type pineconeClient struct {
	config     Config
	httpClient *http.Client
	apiURL     string
}

// PineconeUpsertRequest represents a request to upsert vectors to Pinecone
type PineconeUpsertRequest struct {
	Vectors   []PineconeVector `json:"vectors"`
	Namespace string           `json:"namespace,omitempty"`
}

// PineconeVector represents a vector in Pinecone
type PineconeVector struct {
	ID       string                 `json:"id"`
	Values   []float32              `json:"values"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// PineconeQueryRequest represents a query request to Pinecone
type PineconeQueryRequest struct {
	Vector    []float32 `json:"vector"`
	TopK      int       `json:"topK"`
	Namespace string    `json:"namespace,omitempty"`
}

// PineconeQueryResponse represents a response from a Pinecone query
type PineconeQueryResponse struct {
	Matches []PineconeMatch `json:"matches"`
}

// PineconeMatch represents a match from a Pinecone query
type PineconeMatch struct {
	ID       string                 `json:"id"`
	Score    float32                `json:"score"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// PineconeDeleteRequest represents a delete request to Pinecone
type PineconeDeleteRequest struct {
	IDs       []string `json:"ids"`
	Namespace string   `json:"namespace,omitempty"`
}

// newPineconeClient creates a new Pinecone client
func newPineconeClient(config Config) (Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("Pinecone API key is required")
	}
	if config.Environment == "" {
		return nil, fmt.Errorf("Pinecone environment is required")
	}
	if config.IndexName == "" {
		return nil, fmt.Errorf("Pinecone index name is required")
	}

	apiURL := fmt.Sprintf("https://%s-%s.svc.%s.pinecone.io", 
		config.IndexName, config.IndexName, config.Environment)

	return &pineconeClient{
		config:     config,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		apiURL:     apiURL,
	}, nil
}

// Connect implements the Client interface
func (c *pineconeClient) Connect() error {
	// Test the connection by making a simple request
	url := fmt.Sprintf("%s/describe_index_stats", c.apiURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Api-Key", c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error connecting to Pinecone: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Pinecone: %s - %s", resp.Status, string(body))
	}
	
	return nil
}

// Close implements the Client interface
func (c *pineconeClient) Close() error {
	// Nothing to close for HTTP client
	return nil
}

// UpsertVector implements the Client interface
func (c *pineconeClient) UpsertVector(id string, vector []float32, metadata map[string]interface{}) error {
	vectors := []PineconeVector{
		{
			ID:       id,
			Values:   vector,
			Metadata: metadata,
		},
	}
	
	return c.UpsertVectors([]VectorEntry{
		{
			ID:       id,
			Vector:   vector,
			Metadata: metadata,
		},
	})
}

// UpsertVectors implements the Client interface
func (c *pineconeClient) UpsertVectors(vectors []VectorEntry) error {
	// Convert to Pinecone vectors
	pineconeVectors := make([]PineconeVector, len(vectors))
	for i, v := range vectors {
		pineconeVectors[i] = PineconeVector{
			ID:       v.ID,
			Values:   v.Vector,
			Metadata: v.Metadata,
		}
	}
	
	// Create the request
	requestData := PineconeUpsertRequest{
		Vectors:   pineconeVectors,
		Namespace: c.config.Namespace,
	}
	
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}
	
	// Send the request
	url := fmt.Sprintf("%s/vectors/upsert", c.apiURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Api-Key", c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error upserting vectors: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Pinecone: %s - %s", resp.Status, string(body))
	}
	
	return nil
}

// QueryNearest implements the Client interface
func (c *pineconeClient) QueryNearest(vector []float32, limit int) ([]QueryResult, error) {
	// Create the request
	requestData := PineconeQueryRequest{
		Vector:    vector,
		TopK:      limit,
		Namespace: c.config.Namespace,
	}
	
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}
	
	// Send the request
	url := fmt.Sprintf("%s/query", c.apiURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Api-Key", c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error querying vectors: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from Pinecone: %s - %s", resp.Status, string(body))
	}
	
	// Parse the response
	var pineconeResponse PineconeQueryResponse
	if err := json.NewDecoder(resp.Body).Decode(&pineconeResponse); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	
	// Convert to QueryResult
	results := make([]QueryResult, len(pineconeResponse.Matches))
	for i, match := range pineconeResponse.Matches {
		results[i] = QueryResult{
			ID:       match.ID,
			Score:    match.Score,
			Metadata: match.Metadata,
		}
	}
	
	return results, nil
}

// DeleteVector implements the Client interface
func (c *pineconeClient) DeleteVector(id string) error {
	return c.DeleteVectors([]string{id})
}

// DeleteVectors implements the Client interface
func (c *pineconeClient) DeleteVectors(ids []string) error {
	// Create the request
	requestData := PineconeDeleteRequest{
		IDs:       ids,
		Namespace: c.config.Namespace,
	}
	
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}
	
	// Send the request
	url := fmt.Sprintf("%s/vectors/delete", c.apiURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Api-Key", c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error deleting vectors: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error from Pinecone: %s - %s", resp.Status, string(body))
	}
	
	return nil
}