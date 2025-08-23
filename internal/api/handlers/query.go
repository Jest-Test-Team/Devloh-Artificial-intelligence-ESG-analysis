package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"devlog-ai/pkg/apimodels"
)

// QueryHandler processes natural language queries against logs
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the request
	var query apimodels.QueryRequest
	if err := json.Unmarshal(body, &query); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate the query
	if query.Query == "" {
		http.Error(w, "Query cannot be empty", http.StatusBadRequest)
		return
	}

	// TODO: Implement actual RAG query processing
	// For now, return a mock response
	response := apimodels.QueryResponse{
		Answer: "This is a placeholder response to your query: " + query.Query,
		RelatedLogs: []string{
			"2023-06-01T12:00:00Z ERROR: Connection failed",
			"2023-06-01T12:01:00Z WARN: Retry attempt 1",
		},
		Sources: []string{
			"app-service-1",
			"database-service",
		},
		ProcessingMs: 150,
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	
	// Write the response
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
