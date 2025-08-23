package handlers

import (
	"devlog-ai/internal/api/middleware"
	"devlog-ai/internal/rag"
	"devlog-ai/pkg/apimodels"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var ragService *rag.Service

// InitQueryHandler initializes the query handler
func InitQueryHandler() error {
	var err error
	ragService, err = rag.NewService()
	if err != nil {
		return err
	}
	return nil
}

// QueryHandler processes natural language queries against logs
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		// This should not happen if the auth middleware is working correctly,
		// but we'll handle it just in case
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	// If no project ID was provided in the query, use the one from the authenticated user
	if query.ProjectID == "" {
		query.ProjectID = user.ProjectID
	}

	// Check if RAG service is available
	if ragService == nil {
		log.Println("RAG service not initialized, using mock response")
		// Return a mock response
		response := apimodels.QueryResponse{
			Answer: "RAG service is not available. This is a mock response to your query: " + query.Query,
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
		return
	}

	// Use RAG service to process the query
	response, err := ragService.QueryLogs(query)
	if err != nil {
		log.Printf("Error processing query: %v", err)
		http.Error(w, "Error processing query", http.StatusInternalServerError)
		return
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