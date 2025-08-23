package ingestor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"devlog-ai/internal/api/models"
	"devlog-ai/internal/privacy"
)

// IngestorService handles log ingestion and processing
type IngestorService struct {
	piiDetector *privacy.PIIDetector
}

// NewIngestorService creates a new ingestor service instance
func NewIngestorService() *IngestorService {
	return &IngestorService{
		piiDetector: privacy.NewPIIDetector(),
	}
}

// IngestLog processes and ingests a single log entry
func (i *IngestorService) IngestLog(ctx context.Context, logEntry models.LogEntry) error {
	// Anonymize PII in the log entry
	_ = i.anonymizeLogEntry(logEntry)

	// TODO: Store the anonymized log in the database
	// TODO: Send to Kafka for further processing

	return nil
}

// IngestLogBatch processes and ingests multiple log entries
func (i *IngestorService) IngestLogBatch(ctx context.Context, logs []models.LogEntry) error {
	for _, log := range logs {
		if err := i.IngestLog(ctx, log); err != nil {
			return fmt.Errorf("failed to ingest log %s: %w", log.ID, err)
		}
	}
	return nil
}

// anonymizeLogEntry anonymizes PII in log entry fields
func (i *IngestorService) anonymizeLogEntry(logEntry models.LogEntry) models.LogEntry {
	anonymized := logEntry

	// Anonymize message if it contains PII
	if i.piiDetector.IsPII(logEntry.Message) {
		anonymized.Message = i.piiDetector.DetectAndAnonymize(logEntry.Message)
	}

	// Anonymize metadata if it contains PII
	if logEntry.Metadata != nil {
		anonymized.Metadata = i.piiDetector.AnonymizeLogEntry(logEntry.Metadata)
	}

	return anonymized
}

// ProcessLogFromKafka processes a log message received from Kafka
func (i *IngestorService) ProcessLogFromKafka(ctx context.Context, message []byte) error {
	var logEntry models.LogEntry
	if err := json.Unmarshal(message, &logEntry); err != nil {
		return fmt.Errorf("failed to unmarshal log entry: %w", err)
	}

	return i.IngestLog(ctx, logEntry)
}

// GetIngestionStats returns statistics about log ingestion
func (i *IngestorService) GetIngestionStats(ctx context.Context) (map[string]interface{}, error) {
	// TODO: Implement actual statistics collection
	stats := map[string]interface{}{
		"total_logs_processed": 0,
		"logs_processed_today": 0,
		"pii_detected_count":   0,
		"last_processed_at":    time.Now(),
	}

	return stats, nil
}
