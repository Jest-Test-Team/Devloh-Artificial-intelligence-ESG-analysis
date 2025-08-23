package vectordb

import (
	"errors"
	"fmt"
	"os"
)

// ClientType represents the type of vector database client
type ClientType string

const (
	// PineconeClient is the client type for Pinecone
	PineconeClient ClientType = "pinecone"
	// WeaviateClient is the client type for Weaviate
	WeaviateClient ClientType = "weaviate"
)

// Config represents the configuration for a vector database client
type Config struct {
	Type        ClientType
	APIKey      string
	Environment string // For Pinecone
	URL         string // For Weaviate
	IndexName   string
	Namespace   string // Optional, for Pinecone
	BatchSize   int    // For batch operations
	Dimensions  int    // Vector dimensions
}

// Client is the interface for vector database operations
type Client interface {
	// Connect establishes a connection to the vector database
	Connect() error
	
	// Close closes the connection to the vector database
	Close() error
	
	// UpsertVector inserts or updates a vector in the database
	UpsertVector(id string, vector []float32, metadata map[string]interface{}) error
	
	// UpsertVectors inserts or updates multiple vectors in the database
	UpsertVectors(vectors []VectorEntry) error
	
	// QueryNearest queries for the nearest neighbors of a vector
	QueryNearest(vector []float32, limit int) ([]QueryResult, error)
	
	// DeleteVector deletes a vector from the database
	DeleteVector(id string) error
	
	// DeleteVectors deletes multiple vectors from the database
	DeleteVectors(ids []string) error
}

// VectorEntry represents a vector to be inserted into the database
type VectorEntry struct {
	ID       string
	Vector   []float32
	Metadata map[string]interface{}
}

// QueryResult represents a result from a vector query
type QueryResult struct {
	ID       string
	Score    float32
	Metadata map[string]interface{}
}

// NewClient creates a new vector database client
func NewClient(config Config) (Client, error) {
	switch config.Type {
	case PineconeClient:
		return newPineconeClient(config)
	case WeaviateClient:
		return newWeaviateClient(config)
	default:
		return nil, errors.New("unsupported vector database client type")
	}
}

// GetConfigFromEnv creates a vector database config from environment variables
func GetConfigFromEnv() (Config, error) {
	var config Config
	
	// Get the client type
	clientTypeStr := os.Getenv("VECTORDB_TYPE")
	if clientTypeStr == "" {
		return config, errors.New("VECTORDB_TYPE environment variable not set")
	}
	
	config.Type = ClientType(clientTypeStr)
	
	// Get the common config
	config.IndexName = os.Getenv("VECTORDB_INDEX_NAME")
	if config.IndexName == "" {
		config.IndexName = "devlogai-logs"
	}
	
	config.BatchSize = 100 // Default batch size
	config.Dimensions = 1536 // Default dimensions for OpenAI embeddings
	
	// Type-specific config
	switch config.Type {
	case PineconeClient:
		config.APIKey = os.Getenv("PINECONE_API_KEY")
		if config.APIKey == "" {
			return config, errors.New("PINECONE_API_KEY environment variable not set")
		}
		
		config.Environment = os.Getenv("PINECONE_ENVIRONMENT")
		if config.Environment == "" {
			return config, errors.New("PINECONE_ENVIRONMENT environment variable not set")
		}
		
		config.Namespace = os.Getenv("PINECONE_NAMESPACE")
	
	case WeaviateClient:
		config.APIKey = os.Getenv("WEAVIATE_API_KEY")
		if config.APIKey == "" {
			return config, errors.New("WEAVIATE_API_KEY environment variable not set")
		}
		
		config.URL = os.Getenv("WEAVIATE_URL")
		if config.URL == "" {
			return config, errors.New("WEAVIATE_URL environment variable not set")
		}
	
	default:
		return config, fmt.Errorf("unsupported vector database type: %s", config.Type)
	}
	
	return config, nil
}
