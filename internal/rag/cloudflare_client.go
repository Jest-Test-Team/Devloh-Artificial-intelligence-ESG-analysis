package rag

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// CloudflareClient represents a client for Cloudflare's RAG API
type CloudflareClient struct {
	apiKey     string
	accountID  string
	workerName string
	httpClient *http.Client
}

// CloudflareRAGRequest represents a request to the Cloudflare RAG API
type CloudflareRAGRequest struct {
	Query     string                 `json:"query"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
	ProjectID string                 `json:"project_id,omitempty"`
}

// CloudflareRAGResponse represents a response from the Cloudflare RAG API
type CloudflareRAGResponse struct {
	Answer       string                   `json:"answer"`
	Citations    []CloudflareCitation     `json:"citations,omitempty"`
	RelatedDocs  []CloudflareRelatedDoc   `json:"related_docs,omitempty"`
	ProcessingMs int                      `json:"processing_ms"`
	Metadata     map[string]interface{}   `json:"metadata,omitempty"`
}

// CloudflareCitation represents a citation in a Cloudflare RAG response
type CloudflareCitation struct {
	Text      string                 `json:"text"`
	DocID     string                 `json:"doc_id"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Relevance float64                `json:"relevance,omitempty"`
}

// CloudflareRelatedDoc represents a related document in a Cloudflare RAG response
type CloudflareRelatedDoc struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Score    float64                `json:"score,omitempty"`
}

// NewCloudflareClient creates a new Cloudflare RAG client
func NewCloudflareClient() (*CloudflareClient, error) {
	apiKey := os.Getenv("CLOUDFLARE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("CLOUDFLARE_API_KEY environment variable not set")
	}

	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	if accountID == "" {
		return nil, fmt.Errorf("CLOUDFLARE_ACCOUNT_ID environment variable not set")
	}

	workerName := os.Getenv("CLOUDFLARE_WORKER_NAME")
	if workerName == "" {
		workerName = "devlog-ai-rag"
	}

	return &CloudflareClient{
		apiKey:     apiKey,
		accountID:  accountID,
		workerName: workerName,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// QueryRAG sends a query to the Cloudflare RAG API
func (c *CloudflareClient) QueryRAG(query string, projectID string, options map[string]interface{}) (*CloudflareRAGResponse, error) {
	// Create the request
	requestData := CloudflareRAGRequest{
		Query:     query,
		ProjectID: projectID,
		Options:   options,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Send the request
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/workers/services/%s/rag", c.accountID, c.workerName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from Cloudflare: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	var response CloudflareRAGResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &response, nil
}

// QueryRAGWithContext sends a query to the Cloudflare RAG API with additional context
func (c *CloudflareClient) QueryRAGWithContext(query string, projectID string, context map[string]interface{}, options map[string]interface{}) (*CloudflareRAGResponse, error) {
	// Create the request
	requestData := CloudflareRAGRequest{
		Query:     query,
		ProjectID: projectID,
		Context:   context,
		Options:   options,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Send the request
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/workers/services/%s/rag", c.accountID, c.workerName)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from Cloudflare: %s - %s", resp.Status, string(body))
	}

	// Parse the response
	var response CloudflareRAGResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &response, nil
}
