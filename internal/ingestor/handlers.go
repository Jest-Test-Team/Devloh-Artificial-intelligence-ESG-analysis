package ingestor

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"devlog-ai/internal/api/models"
	"devlog-ai/internal/api/middleware" // 修正導入路徑
)

// LogHandler handles incoming logs
func (s *Service) LogHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context if authentication is enabled
	user, err := middleware.GetUserFromContext(r.Context())
	if err != nil {
		// If authentication is enabled, this would be an error
		// For now, we'll just log it and continue
		log.Printf("Warning: No authenticated user found: %v", err)
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
		// If authentication is enabled and user is available, use their project ID
		if user != nil {
			logBatch.ProjectID = user.ProjectID
		} else {
			http.Error(w, "ProjectID is required", http.StatusBadRequest)
			return
		}
	}
	
	// If authentication is enabled, ensure user can only send logs to their own project
	if user != nil && logBatch.ProjectID != user.ProjectID {
		http.Error(w, "Unauthorized: can only send logs to your own project", http.StatusForbidden)
		return
	}

	if len(logBatch.Logs) == 0 {
		http.Error(w, "Log batch cannot be empty", http.StatusBadRequest)
		return
	}

	// Process the log batch
	if err := s.ProcessLogBatch(logBatch); err != nil {
		log.Printf("Error processing log batch: %v", err)
		http.Error(w, "Error processing logs", http.StatusInternalServerError)
		return
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status":"accepted"}`))
}