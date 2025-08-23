package ingestor

import (
	"encoding/json"
	"net/http"

	"devlog-ai/internal/api/models"
)

// IngestorHandler handles HTTP requests for the ingestor service
type IngestorHandler struct {
	service *IngestorService
}

// NewIngestorHandler creates a new ingestor handler
func NewIngestorHandler(service *IngestorService) *IngestorHandler {
	return &IngestorHandler{
		service: service,
	}
}

// IngestLogHandler handles log ingestion requests
func (h *IngestorHandler) IngestLogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logEntry models.LogEntry
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.IngestLog(r.Context(), logEntry); err != nil {
		http.Error(w, "Failed to ingest log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Log ingested successfully",
	})
}

// IngestBatchHandler handles batch log ingestion requests
func (h *IngestorHandler) IngestBatchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logBatch models.LogBatch
	if err := json.NewDecoder(r.Body).Decode(&logBatch); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.IngestLogBatch(r.Context(), logBatch.Logs); err != nil {
		http.Error(w, "Failed to ingest log batch", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Log batch ingested successfully",
	})
}

// HealthHandler handles health check requests
func (h *IngestorHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "ingestor",
	})
}
