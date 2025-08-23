package models

import "time"

// LogEntry represents a single log entry in the system
type LogEntry struct {
	ID          string            `json:"id"`
	Timestamp   time.Time         `json:"timestamp"`
	Level       string            `json:"level"`
	Message     string            `json:"message"`
	Source      string            `json:"source"`
	ServiceName string            `json:"service_name"`
	ProjectID   string            `json:"project_id"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// LogBatch represents a batch of logs sent to the system
type LogBatch struct {
	ProjectID string     `json:"project_id"`
	Logs      []LogEntry `json:"logs"`
}
