package processor

import (
	"context"
	"fmt"
	"time"

	"devlog-ai/internal/api/models"
	"devlog-ai/internal/embedding"
	"devlog-ai/internal/vectordb"
)

// ProcessorService handles log processing and analysis
type ProcessorService struct {
	embeddingClient embedding.EmbeddingClient
	vectorDBClient  vectordb.VectorDBClient
}

// NewProcessorService creates a new processor service instance
func NewProcessorService(embeddingClient embedding.EmbeddingClient, vectorDBClient vectordb.VectorDBClient) *ProcessorService {
	return &ProcessorService{
		embeddingClient: embeddingClient,
		vectorDBClient:  vectorDBClient,
	}
}

// ProcessLog processes a log entry and stores it in the vector database
func (p *ProcessorService) ProcessLog(ctx context.Context, logEntry models.LogEntry) error {
	// Create embedding for the log message
	embedding, err := p.embeddingClient.CreateEmbedding(logEntry.Message)
	if err != nil {
		return fmt.Errorf("failed to create embedding: %w", err)
	}

	// Prepare vector for storage
	vector := vectordb.Vector{
		ID: logEntry.ID,
		Vector: embedding,
		Metadata: map[string]interface{}{
			"service":   logEntry.Service,
			"level":     logEntry.Level,
			"timestamp": logEntry.Timestamp,
			"user_id":   logEntry.UserID,
			"metadata":  logEntry.Metadata,
		},
	}

	// Store in vector database
	if err := p.vectorDBClient.Upsert(ctx, "logs", []vectordb.Vector{vector}); err != nil {
		return fmt.Errorf("failed to store vector: %w", err)
	}

	return nil
}

// ProcessLogBatch processes multiple log entries
func (p *ProcessorService) ProcessLogBatch(ctx context.Context, logs []models.LogEntry) error {
	for _, log := range logs {
		if err := p.ProcessLog(ctx, log); err != nil {
			return fmt.Errorf("failed to process log %s: %w", log.ID, err)
		}
	}
	return nil
}

// SearchSimilarLogs finds logs similar to the given query
func (p *ProcessorService) SearchSimilarLogs(ctx context.Context, query string, limit int) ([]models.LogEntry, error) {
	// Create embedding for the query
	queryEmbedding, err := p.embeddingClient.CreateEmbedding(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create query embedding: %w", err)
	}

	// Search for similar vectors
	searchVector := vectordb.Vector{
		Vector: queryEmbedding,
	}

	results, err := p.vectorDBClient.Search(ctx, "logs", searchVector, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector database: %w", err)
	}

	// Convert results back to log entries
	var logEntries []models.LogEntry
	for _, result := range results {
		if metadata, ok := result.Metadata["service"].(string); ok {
			logEntry := models.LogEntry{
				ID:        result.ID,
				Service:   metadata,
				Level:     getStringFromMetadata(result.Metadata, "level"),
				Message:   getStringFromMetadata(result.Metadata, "message"),
				Timestamp: getStringFromMetadata(result.Metadata, "timestamp"),
				UserID:    getStringFromMetadata(result.Metadata, "user_id"),
			}
			logEntries = append(logEntries, logEntry)
		}
	}

	return logEntries, nil
}

// getStringFromMetadata safely extracts a string value from metadata
func getStringFromMetadata(metadata map[string]interface{}, key string) string {
	if value, ok := metadata[key].(string); ok {
		return value
	}
	return ""
}

// GetProcessingStats returns statistics about log processing
func (p *ProcessorService) GetProcessingStats(ctx context.Context) (map[string]interface{}, error) {
	// TODO: Implement actual statistics collection
	stats := map[string]interface{}{
		"total_logs_processed": 0,
		"logs_processed_today": 0,
		"last_processed_at":    time.Now(),
	}

	return stats, nil
}
