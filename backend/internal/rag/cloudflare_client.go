package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CloudflareClient provides access to Cloudflare AI services
type CloudflareClient struct {
	apiToken string
	baseURL  string
	client   *http.Client
}

// CloudflareAIRequest represents a request to Cloudflare AI
type CloudflareAIRequest struct {
	Prompt      string                 `json:"prompt"`
	Model       string                 `json:"model"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

// CloudflareAIResponse represents a response from Cloudflare AI
type CloudflareAIResponse struct {
	Response string `json:"response"`
	Model    string `json:"model"`
	Usage    struct {
		PromptTokens   int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens    int `json:"total_tokens"`
	} `json:"usage"`
}

// NewCloudflareClient creates a new Cloudflare client
func NewCloudflareClient(apiToken string) *CloudflareClient {
	return &CloudflareClient{
		apiToken: apiToken,
		baseURL:  "https://api.cloudflare.com/client/v4/ai/run",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GenerateText generates text using Cloudflare AI
func (c *CloudflareClient) GenerateText(ctx context.Context, prompt string, model string) (string, error) {
	request := CloudflareAIRequest{
		Prompt:      prompt,
		Model:       model,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response CloudflareAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Response, nil
}

// Close closes the Cloudflare client
func (c *CloudflareClient) Close() error {
	// No resources to close for HTTP client
	return nil
}
