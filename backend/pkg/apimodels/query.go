package apimodels

// QueryRequest represents a natural language query from a user
type QueryRequest struct {
	Query     string                 `json:"query"`
	ProjectID string                 `json:"project_id,omitempty"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

// QueryResponse represents the response to a user query
type QueryResponse struct {
	Answer        string                     `json:"answer"`
	RelevantLogs  []map[string]interface{}  `json:"relevant_logs,omitempty"`
	RelatedLogs   []string                  `json:"related_logs,omitempty"`
	Sources       []string                  `json:"sources,omitempty"`
	ProcessingMs  int64                     `json:"processing_ms"`
}
