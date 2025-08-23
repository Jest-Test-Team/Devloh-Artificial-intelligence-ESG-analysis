package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	// 創建一個簡單的處理函數
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// 創建CORS中間件
	corsMiddleware := CORSMiddleware(nextHandler)

	// 測試正常請求
	t.Run("Normal Request", func(t *testing.T) {
		// 創建請求
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		// 創建響應記錄器
		rr := httptest.NewRecorder()

		// 調用中間件
		corsMiddleware.ServeHTTP(rr, req)

		// 檢查狀態碼
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// 檢查CORS頭
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("expected Access-Control-Allow-Origin header to be '*', got %v",
				rr.Header().Get("Access-Control-Allow-Origin"))
		}

		if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Errorf("unexpected Access-Control-Allow-Methods header: %v",
				rr.Header().Get("Access-Control-Allow-Methods"))
		}

		if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
			t.Errorf("unexpected Access-Control-Allow-Headers header: %v",
				rr.Header().Get("Access-Control-Allow-Headers"))
		}
	})

	// 測試預檢請求
	t.Run("Preflight Request", func(t *testing.T) {
		// 創建OPTIONS請求
		req, err := http.NewRequest("OPTIONS", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		// 創建響應記錄器
		rr := httptest.NewRecorder()

		// 調用中間件
		corsMiddleware.ServeHTTP(rr, req)

		// 檢查狀態碼（應該是200 OK）
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// 檢查CORS頭
		if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Errorf("expected Access-Control-Allow-Origin header to be '*', got %v",
				rr.Header().Get("Access-Control-Allow-Origin"))
		}

		if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
			t.Errorf("unexpected Access-Control-Allow-Methods header: %v",
				rr.Header().Get("Access-Control-Allow-Methods"))
		}

		if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
			t.Errorf("unexpected Access-Control-Allow-Headers header: %v",
				rr.Header().Get("Access-Control-Allow-Headers"))
		}
	})
}
