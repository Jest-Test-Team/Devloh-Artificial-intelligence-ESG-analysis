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

// 模擬RAG服務
type mockRAGService struct{}

func (m *mockRAGService) Query(query string, options map[string]interface{}, userID string) (*apimodels.QueryResponse, error) {
	// 返回模擬數據
	return &apimodels.QueryResponse{
		Answer: "這是對查詢的模擬回答",
		RelevantLogs: []map[string]interface{}{
			{
				"timestamp": "2023-07-15T12:00:00Z",
				"service":   "api-service",
				"level":     "INFO",
				"message":   "這是一個相關的日誌條目",
			},
		},
	}, nil
}

func TestQueryHandler(t *testing.T) {
	// 創建模擬RAG服務
	mockService := &mockRAGService{}

	// 創建查詢處理程序
	handler := QueryHandler(mockService)

	// 創建請求體
	requestBody := apimodels.QueryRequest{
		Query: "查詢系統錯誤",
		Options: map[string]interface{}{
			"max_results":     5,
			"time_range_hours": 24,
		},
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// 創建一個新的請求
	req, err := http.NewRequest("POST", "/query", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 模擬認證中間件
	// 將用戶ID添加到請求上下文
	ctx := req.Context()
	ctx = middleware.WithUserID(ctx, "test-user-123")
	req = req.WithContext(ctx)

	// 創建一個ResponseRecorder來記錄響應
	rr := httptest.NewRecorder()

	// 直接調用處理函數
	handler.ServeHTTP(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 解析響應體
	var response apimodels.QueryResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	// 檢查響應內容
	if response.Answer != "這是對查詢的模擬回答" {
		t.Errorf("handler returned unexpected answer: got %v want %v",
			response.Answer, "這是對查詢的模擬回答")
	}

	if len(response.RelevantLogs) != 1 {
		t.Errorf("handler returned unexpected number of logs: got %v want %v",
			len(response.RelevantLogs), 1)
	}
}

// 測試無效請求
func TestQueryHandlerInvalidRequest(t *testing.T) {
	// 創建模擬RAG服務
	mockService := &mockRAGService{}

	// 創建查詢處理程序
	handler := QueryHandler(mockService)

	// 創建無效的請求體（缺少查詢字段）
	requestBody := map[string]interface{}{
		"options": map[string]interface{}{
			"max_results": 5,
		},
	}
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	// 創建一個新的請求
	req, err := http.NewRequest("POST", "/query", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 模擬認證中間件
	ctx := req.Context()
	ctx = middleware.WithUserID(ctx, "test-user-123")
	req = req.WithContext(ctx)

	// 創建一個ResponseRecorder來記錄響應
	rr := httptest.NewRecorder()

	// 直接調用處理函數
	handler.ServeHTTP(rr, req)

	// 檢查狀態碼（應該是400 Bad Request）
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
