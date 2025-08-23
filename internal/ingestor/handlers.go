package ingestor

import (
	"encoding/json"
	"io"
	"net/http"

	"devlog-ai/internal/api/models"
)

// LogHandler handles incoming logs
func (s *Service) LogHandler(w http.ResponseWriter, r *http.Request) {
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

	// Parse the log batch
	var logBatch models.LogBatch
	if err := json.Unmarshal(body, &logBatch); err != nil {
		http.Error(w, "Invalid log format", http.StatusBadRequest)
		return
	}

	// Validate the log batch
	if logBatch.ProjectID == "" {
		http.Error(w, "ProjectID is required", http.StatusBadRequest)
		return
	}
	if len(logBatch.Logs) == 0 {
		http.Error(w, "Log batch cannot be empty", http.StatusBadRequest)
		return
	}

	// Process the log batch
	if err := s.ProcessLogBatch(logBatch); err != nil {
		http.Error(w, "Error processing logs", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}
