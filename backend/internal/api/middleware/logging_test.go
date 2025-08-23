package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// 捕獲日誌輸出
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(nil) // 恢復默認輸出
	}()

	// 創建一個簡單的處理函數
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 創建日誌中間件
	loggingMiddleware := LoggingMiddleware(nextHandler)

	// 創建請求
	req, err := http.NewRequest("GET", "/test?param=value", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "Test-Agent")
	req.Header.Set("Referer", "http://example.com")
	req.RemoteAddr = "192.168.1.1:12345"

	// 創建響應記錄器
	rr := httptest.NewRecorder()

	// 調用中間件
	loggingMiddleware.ServeHTTP(rr, req)

	// 檢查狀態碼
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 檢查日誌輸出
	logOutput := buf.String()
	
	// 檢查日誌是否包含關鍵信息
	expectedParts := []string{
		"GET",
		"/test?param=value",
		"192.168.1.1",
		"Test-Agent",
		"http://example.com",
		"200",
	}

	for _, part := range expectedParts {
		if !strings.Contains(logOutput, part) {
			t.Errorf("log output missing expected part: %s", part)
		}
	}
}
