package middleware

import (
	"context"
	"net/http"
)

// UserIDKey is the context key for user ID
type UserIDKey string

const (
	UserIDContextKey UserIDKey = "user_id"
)

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDContextKey, userID)
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDContextKey).(string); ok {
		return userID
	}
	return ""
}

// AuthMiddleware authenticates requests
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now, just pass through - implement actual auth logic later
		// TODO: Implement JWT token validation, API key validation, etc.
		next.ServeHTTP(w, r)
	})
}

// RequireAuth middleware that requires authentication
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now, just pass through - implement actual auth logic later
		// TODO: Check for valid authentication and return 401 if not authenticated
		next.ServeHTTP(w, r)
	})
}
