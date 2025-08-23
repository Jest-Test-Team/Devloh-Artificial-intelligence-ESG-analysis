package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"devlog-ai/internal/rag"
	"devlog-ai/pkg/apimodels"
)

// QueryHandler handles natural language query requests
type QueryHandler struct {
	ragService *rag.RAGService
}

// NewQueryHandler creates a new query handler
func NewQueryHandler(ragService *rag.RAGService) *QueryHandler {
	return &QueryHandler{
		ragService: ragService,
	}
}

// QueryHandler handles natural language queries
func (h *QueryHandler) QueryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request apimodels.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	// Process the query using RAG service
	answer, err := h.ragService.Query(r.Context(), request.Query, "logs")
	if err != nil {
		http.Error(w, "Failed to process query", http.StatusInternalServerError)
		return
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Create response
	response := apimodels.QueryResponse{
		Answer:       answer,
		RelatedLogs:  []string{}, // TODO: Implement related logs retrieval
		Sources:      []string{}, // TODO: Implement sources retrieval
		ProcessingMs: processingTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HealthHandler handles health check requests
func (h *QueryHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "query",
	})
}
