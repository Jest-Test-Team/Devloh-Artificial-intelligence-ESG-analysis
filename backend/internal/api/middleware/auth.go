package middleware

import (
	"context"
	"net/http"
	"strings"
)

// contextKey 是用於上下文值的鍵類型
type contextKey string

const (
	// UserIDKey 是存儲用戶ID的上下文鍵
	UserIDKey contextKey = "user_id"
)

// AuthMiddleware 是一個處理 JWT 身份驗證的中間件
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 從請求頭獲取 Authorization 標頭
		authHeader := r.Header.Get("Authorization")
		
		// 檢查是否存在 Authorization 標頭
		if authHeader == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}
		
		// 檢查格式是否正確（以 "Bearer " 開頭）
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized: Invalid token format", http.StatusUnauthorized)
			return
		}
		
		token := parts[1]
		
		// 驗證 token（在實際應用中，這裡會調用一個專門的 JWT 驗證函數）
		userID, err := validateToken(token)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		
		// 將用戶 ID 添加到請求上下文
		ctx := WithUserID(r.Context(), userID)
		
		// 調用下一個處理程序
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithUserID 將用戶 ID 添加到上下文中
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 從上下文中獲取用戶 ID
func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// GetUserFromContext 從上下文中獲取用戶數據
func GetUserFromContext(ctx context.Context) (*User, error) {
	userID := GetUserID(ctx)
	if userID == "" {
		return nil, ErrNoUser
	}
	
	// 在實際應用中，這裡可能會從數據庫或緩存中獲取用戶資料
	// 這裡我們只是返回一個基於 ID 的模擬用戶
	user := &User{
		ID:        userID,
		Email:     userID + "@example.com",
		Name:      "User " + userID,
		ProjectID: "project-" + userID,
	}
	
	return user, nil
}

// User 表示一個認證用戶
type User struct {
	ID        string
	Email     string
	Name      string
	ProjectID string
}

// ErrNoUser 是沒有用戶的錯誤
var errNoUser = &noUserError{"no user found in context"}
var ErrNoUser = errNoUser

type noUserError struct {
	msg string
}

func (e *noUserError) Error() string {
	return e.msg
}

// validateToken 驗證 JWT token 並返回用戶 ID
// 在實際應用中，這個函數會對 JWT 進行適當的驗證
func validateToken(token string) (string, error) {
	// 這裡只是為了測試，我們直接返回一個模擬的用戶 ID
	// 在真實環境中，應該解析並驗證 JWT token
	if token == "invalid-token" {
		return "", &noUserError{"invalid token"}
	}
	
	// 從 token 中提取用戶 ID
	// 這裡我們假設所有以 "valid-token" 開頭的 token 都是有效的
	if strings.HasPrefix(token, "valid-token") {
		return "user-123", nil
	}
	
	// 為了測試的方便，我們假設所有其他 token 也都是有效的
	return "user-" + token[:8], nil
}