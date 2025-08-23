package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	// 創建一個新的請求
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 創建一個ResponseRecorder來記錄響應
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	// 直接調用處理函數
	handler.ServeHTTP(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢查響應體
	expected := `{"status":"ok"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
