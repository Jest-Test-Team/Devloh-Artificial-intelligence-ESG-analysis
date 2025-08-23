package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIngestLogsHandler(t *testing.T) {
	// Initialize log handlers for testing
	if err := InitLogHandlers(); err != nil {
		t.Fatalf("Failed to initialize log handlers for testing: %v", err)
	}

	// Test cases
	testCases := []struct {
		name           string
		method         string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Valid Log Batch",
			method: "POST",
			requestBody: map[string]interface{}{
				"project_id": "test-project",
				"logs": []map[string]interface{}{
					{
						"timestamp":    time.Now().Format(time.RFC3339),
						"level":        "ERROR",
						"message":      "Test error message",
						"source":       "test-source",
						"service_name": "test-service",
					},
				},
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Method Not Allowed",
			method:         "GET",
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "Missing Project ID",
			method: "POST",
			requestBody: map[string]interface{}{
				"logs": []map[string]interface{}{
					{
						"timestamp":    time.Now().Format(time.RFC3339),
						"level":        "ERROR",
						"message":      "Test error message",
						"source":       "test-source",
						"service_name": "test-service",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "Empty Logs Array",
			method: "POST",
			requestBody: map[string]interface{}{
				"project_id": "test-project",
				"logs":       []map[string]interface{}{},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			var requestBody []byte
			var err error
			if tc.requestBody != nil {
				requestBody, err = json.Marshal(tc.requestBody)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Create a request
			req, err := http.NewRequest(tc.method, "/api/v1/logs", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Add content type header for POST requests
			if tc.method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}

			// Create a ResponseRecorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(IngestLogsHandler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// For successful requests, verify response
			if tc.expectedStatus == http.StatusAccepted {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("failed to parse response: %v", err)
				}

				// Check status field
				if status, exists := response["status"]; !exists || status != "accepted" {
					t.Errorf("response missing or incorrect 'status' field: %v", response)
				}
			}
		})
	}
}
