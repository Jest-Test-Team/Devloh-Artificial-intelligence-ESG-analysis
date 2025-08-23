package embedding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ClientType represents the type of embedding client
type ClientType string

const (
	// OpenAIClient is the client type for OpenAI
	OpenAIClient ClientType = "openai"
	// GeminiClient is the client type for Google Gemini
	GeminiClient ClientType = "gemini"
)

// Config represents the configuration for an embedding client
type Config struct {
	Type      ClientType
	APIKey    string
	ModelName string
	BatchSize int // For batch operations
}

// Client is the interface for embedding operations
type Client interface {
	// EmbedText creates an embedding for a single text
	EmbedText(text string) ([]float32, error)
	
	// EmbedTexts creates embeddings for multiple texts
	EmbedTexts(texts []string) ([][]float32, error)
}

// NewClient creates a new embedding client
func NewClient(config Config) (Client, error) {
	switch config.Type {
	case OpenAIClient:
		return newOpenAIClient(config)
	case GeminiClient:
		return newGeminiClient(config)
	default:
		return nil, errors.New("unsupported embedding client type")
	}
}

// GetConfigFromEnv creates an embedding config from environment variables
func GetConfigFromEnv() (Config, error) {
	var config Config
	
	// Get the client type
	clientTypeStr := os.Getenv("EMBEDDING_TYPE")
	if clientTypeStr == "" {
		clientTypeStr = string(OpenAIClient) // Default to OpenAI
	}
	
	config.Type = ClientType(clientTypeStr)
	
	// Type-specific config
	switch config.Type {
	case OpenAIClient:
		config.APIKey = os.Getenv("OPENAI_API_KEY")
		if config.APIKey == "" {
			return config, errors.New("OPENAI_API_KEY environment variable not set")
		}
		
		config.ModelName = os.Getenv("OPENAI_EMBEDDING_MODEL")
		if config.ModelName == "" {
			config.ModelName = "text-embedding-ada-002" // Default model
		}
		
	case GeminiClient:
		config.APIKey = os.Getenv("GEMINI_API_KEY")
		if config.APIKey == "" {
			return config, errors.New("GEMINI_API_KEY environment variable not set")
		}
		
		config.ModelName = os.Getenv("GEMINI_EMBEDDING_MODEL")
		if config.ModelName == "" {
			config.ModelName = "embedding-001" // Default model
		}
		
	default:
		return config, fmt.Errorf("unsupported embedding client type: %s", config.Type)
	}
	
	batchSizeStr := os.Getenv("EMBEDDING_BATCH_SIZE")
	if batchSizeStr != "" {
		var batchSize int
		if _, err := fmt.Sscanf(batchSizeStr, "%d", &batchSize); err == nil && batchSize > 0 {
			config.BatchSize = batchSize
		}
	}
	
	if config.BatchSize <= 0 {
		config.BatchSize = 10 // Default batch size
	}
	
	return config, nil
}

// OpenAIClient is a client for the OpenAI embedding API
type openAIClient struct {
	config     Config
	httpClient *http.Client
}

// OpenAIEmbeddingRequest represents a request to the OpenAI embedding API
type openAIEmbeddingRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

// OpenAIEmbeddingResponse represents a response from the OpenAI embedding API
type openAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

// newOpenAIClient creates a new OpenAI client
func newOpenAIClient(config Config) (*openAIClient, error) {
	if config.APIKey == "" {
		return nil, errors.New("OpenAI API key is required")
	}
	
	return &openAIClient{
		config:     config,
		httpClient: &http.Client{},
	}, nil
}

// EmbedText implements the Client interface
func (c *openAIClient) EmbedText(text string) ([]float32, error) {
	embeddings, err := c.EmbedTexts([]string{text})
	if err != nil {
		return nil, err
	}
	
	return embeddings[0], nil
}

// EmbedTexts implements the Client interface
func (c *openAIClient) EmbedTexts(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, errors.New("no texts provided")
	}
	
	// Create the request
	requestData := openAIEmbeddingRequest{
		Input: texts,
		Model: c.config.ModelName,
	}
	
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}
	
	// Send the request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from OpenAI: %s - %s", resp.Status, string(body))
	}
	
	// Parse the response
	var openAIResponse openAIEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResponse); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	
	// Extract the embeddings
	embeddings := make([][]float32, len(openAIResponse.Data))
	for i, data := range openAIResponse.Data {
		embeddings[i] = data.Embedding
	}
	
	return embeddings, nil
}

// GeminiClient is a client for the Google Gemini embedding API
type geminiClient struct {
	config     Config
	httpClient *http.Client
}

// GeminiEmbeddingRequest represents a request to the Gemini embedding API
type geminiEmbeddingRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
	Model string `json:"model"`
}

// GeminiEmbeddingResponse represents a response from the Gemini embedding API
type geminiEmbeddingResponse struct {
	Embeddings []struct {
		Values []float32 `json:"values"`
	} `json:"embeddings"`
}

// newGeminiClient creates a new Gemini client
func newGeminiClient(config Config) (*geminiClient, error) {
	if config.APIKey == "" {
		return nil, errors.New("Gemini API key is required")
	}
	
	return &geminiClient{
		config:     config,
		httpClient: &http.Client{},
	}, nil
}

// EmbedText implements the Client interface
func (c *geminiClient) EmbedText(text string) ([]float32, error) {
	embeddings, err := c.EmbedTexts([]string{text})
	if err != nil {
		return nil, err
	}
	
	return embeddings[0], nil
}

// EmbedTexts implements the Client interface
func (c *geminiClient) EmbedTexts(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, errors.New("no texts provided")
	}
	
	// Create the contents array
	contents := make([]struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}, len(texts))
	
	for i, text := range texts {
		contents[i] = struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			Parts: []struct {
				Text string `json:"text"`
			}{
				{Text: text},
			},
		}
	}
	
	// Create the request
	requestData := geminiEmbeddingRequest{
		Contents: contents,
		Model:    c.config.ModelName,
	}
	
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}
	
	// Send the request
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/%s:embedContent?key=%s", 
		c.config.ModelName, c.config.APIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error from Gemini: %s - %s", resp.Status, string(body))
	}
	
	// Parse the response
	var geminiResponse geminiEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResponse); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	
	// Extract the embeddings
	embeddings := make([][]float32, len(geminiResponse.Embeddings))
	for i, embedding := range geminiResponse.Embeddings {
		embeddings[i] = embedding.Values
	}
	
	return embeddings, nil
}
