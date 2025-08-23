package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"devlog-ai/internal/api/middleware"
	"devlog-ai/pkg/apimodels"
)

// 模擬 RAG 服務
type mockRAGService struct{}

func (m *mockRAGService) Query(query string, options map[string]interface{}) (string, error) {
	return "This is a mock response for query: " + query, nil
}

func TestQueryHandler_HandleQuery(t *testing.T) {
	// 創建模擬 RAG 服務
	mockService := &mockRAGService{}

	// 創建查詢處理程序
	handler := NewQueryHandler(mockService)

	// 創建測試請求
	queryRequest := apimodels.QueryRequest{
		Query: "test query",
		Options: map[string]interface{}{
			"max_results": 5,
		},
	}

	reqBody, _ := json.Marshal(queryRequest)
	req, err := http.NewRequest("POST", "/query", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 設置請求頭
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")

	// 創建響應記錄器
	rr := httptest.NewRecorder()

	// 創建處理程序鏈，包括身份驗證中間件
	authMiddleware := middleware.AuthMiddleware(http.HandlerFunc(handler.HandleQuery))
	authMiddleware.ServeHTTP(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 檢查響應內容
	var response apimodels.QueryResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedPrefix := "This is a mock response for query: test query"
	if response.Answer != expectedPrefix {
		t.Errorf("Handler returned unexpected response: got %v want %v", response.Answer, expectedPrefix)
	}
}

func TestQueryHandler_InvalidMethod(t *testing.T) {
	// 創建模擬 RAG 服務
	mockService := &mockRAGService{}

	// 創建查詢處理程序
	handler := NewQueryHandler(mockService)

	// 創建 GET 請求（應該被拒絕）
	req, err := http.NewRequest("GET", "/query", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 創建響應記錄器
	rr := httptest.NewRecorder()

	// 處理請求
	handler.HandleQuery(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestQueryHandler_InvalidRequestBody(t *testing.T) {
	// 創建模擬 RAG 服務
	mockService := &mockRAGService{}

	// 創建查詢處理程序
	handler := NewQueryHandler(mockService)

	// 創建無效的請求體
	req, err := http.NewRequest("POST", "/query", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 設置請求頭
	req.Header.Set("Content-Type", "application/json")

	// 創建響應記錄器
	rr := httptest.NewRecorder()

	// 處理請求
	handler.HandleQuery(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}