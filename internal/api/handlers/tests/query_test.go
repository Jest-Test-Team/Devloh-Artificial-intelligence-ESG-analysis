package tests

import (
	"bytes"
	"devlog-ai/internal/api/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryHandler(t *testing.T) {
	// Test cases
	testCases := []struct {
		name           string
		method         string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Valid Query",
			method: "POST",
			requestBody: map[string]interface{}{
				"query":      "show me error logs",
				"project_id": "test-project",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Method Not Allowed",
			method:         "GET",
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "Empty Query",
			method: "POST",
			requestBody: map[string]interface{}{
				"query":      "",
				"project_id": "test-project",
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
			req, err := http.NewRequest(tc.method, "/api/v1/query", bytes.NewBuffer(requestBody))
			if err != nil {
				t.Fatal(err)
			}

			// Add content type header for POST requests
			if tc.method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}

			// Create a ResponseRecorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.QueryHandler)

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tc.expectedStatus)
			}

			// For successful requests, verify response structure
			if tc.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("failed to parse response: %v", err)
				}

				// Check required fields
				if _, exists := response["answer"]; !exists {
					t.Errorf("response missing 'answer' field")
				}
				if _, exists := response["processing_ms"]; !exists {
					t.Errorf("response missing 'processing_ms' field")
				}
			}
		})
	}
}
