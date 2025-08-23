package processor

import (
	"encoding/json"
	"net/http"
	"strconv"

	"devlog-ai/internal/api/models"
)

// ProcessorHandler handles HTTP requests for the processor service
type ProcessorHandler struct {
	service *ProcessorService
}

// NewProcessorHandler creates a new processor handler
func NewProcessorHandler(service *ProcessorService) *ProcessorHandler {
	return &ProcessorHandler{
		service: service,
	}
}

// ProcessLogHandler handles log processing requests
func (h *ProcessorHandler) ProcessLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logEntry models.LogEntry
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.ProcessLog(r.Context(), logEntry); err != nil {
		http.Error(w, "Failed to process log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Log processed successfully",
	})
}

// SearchLogsHandler handles log search requests
func (h *ProcessorHandler) SearchLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	logs, err := h.service.SearchSimilarLogs(r.Context(), query, limit)
	if err != nil {
		http.Error(w, "Failed to search logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"logs":   logs,
		"count":  len(logs),
	})
}

// HealthHandler handles health check requests
func (h *ProcessorHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "processor",
	})
}
