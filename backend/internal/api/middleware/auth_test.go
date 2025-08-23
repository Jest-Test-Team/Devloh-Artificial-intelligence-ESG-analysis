package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// 創建一個簡單的處理函數來檢查上下文中的用戶ID
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r.Context())
		if userID == "" {
			t.Error("user ID not found in context")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(userID))
	})

	// 創建認證中間件
	authMiddleware := AuthMiddleware(nextHandler)

	// 測試有效令牌
	t.Run("Valid Token", func(t *testing.T) {
		// 創建請求
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		// 添加有效的認證頭
		req.Header.Set("Authorization", "Bearer valid-token-123")

		// 創建響應記錄器
		rr := httptest.NewRecorder()

		// 調用中間件
		authMiddleware.ServeHTTP(rr, req)

		// 檢查狀態碼
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// 檢查響應體（應該包含用戶ID）
		if rr.Body.String() == "" {
			t.Error("expected non-empty user ID in response")
		}
	})

	// 測試缺少令牌
	t.Run("Missing Token", func(t *testing.T) {
		// 創建請求
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		// 不添加認證頭

		// 創建響應記錄器
		rr := httptest.NewRecorder()

		// 調用中間件
		authMiddleware.ServeHTTP(rr, req)

		// 檢查狀態碼（應該是401 Unauthorized）
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	// 測試無效令牌格式
	t.Run("Invalid Token Format", func(t *testing.T) {
		// 創建請求
		req, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		// 添加無效格式的認證頭
		req.Header.Set("Authorization", "InvalidFormat")

		// 創建響應記錄器
		rr := httptest.NewRecorder()

		// 調用中間件
		authMiddleware.ServeHTTP(rr, req)

		// 檢查狀態碼（應該是401 Unauthorized）
		if status := rr.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})
}

func TestGetUserID(t *testing.T) {
	// 測試從上下文獲取用戶ID
	req, _ := http.NewRequest("GET", "/test", nil)
	ctx := req.Context()
	
	// 測試空上下文
	userID := GetUserID(ctx)
	if userID != "" {
		t.Errorf("expected empty user ID from empty context, got %v", userID)
	}

	// 測試帶有用戶ID的上下文
	ctx = WithUserID(ctx, "test-user-456")
	userID = GetUserID(ctx)
	if userID != "test-user-456" {
		t.Errorf("expected user ID 'test-user-456', got %v", userID)
	}
}
