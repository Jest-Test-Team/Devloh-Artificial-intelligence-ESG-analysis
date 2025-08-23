package models

import (
	"time"
)

// LogEntry represents a log entry in the system
type LogEntry struct {
	ID        string                 `json:"id"`
	Service   string                 `json:"service"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	UserID    string                 `json:"user_id,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// LogBatch represents a batch of log entries
type LogBatch struct {
	Logs []LogEntry `json:"logs"`
}

// LogQuery represents a query for searching logs
type LogQuery struct {
	Service   string    `json:"service,omitempty"`
	Level     string    `json:"level,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}
