package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// User represents authenticated user information
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
}

// contextKey is a custom type to avoid collisions in the context
type contextKey string

// UserContextKey is the key used to store and retrieve the user from the context
const UserContextKey contextKey = "user"

// AuthResponse is the response from an auth service
type AuthResponse struct {
	Valid bool `json:"valid"`
	User  User `json:"user"`
}

// AuthMiddleware checks for a valid JWT token in the Authorization header
func AuthMiddleware(publicPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the path should bypass authentication
			for _, path := range publicPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			// Get the token
			token := parts[1]
			if token == "" {
				http.Error(w, "Empty token", http.StatusUnauthorized)
				return
			}

			// In production, we would verify the token with an auth service
			// For now, we'll validate using a simple mock function
			valid, user, err := validateToken(token)
			if err != nil {
				log.Printf("Error validating token: %v", err)
				http.Error(w, "Internal server error during authentication", http.StatusInternalServerError)
				return
			}

			if !valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Add the user to the context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			
			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the User from the context
func GetUserFromContext(ctx context.Context) (*User, error) {
	user, ok := ctx.Value(UserContextKey).(*User)
	if !ok {
		return nil, fmt.Errorf("no user found in context")
	}
	return user, nil
}

// For development only: mock token validation
func validateToken(token string) (bool, *User, error) {
	// In production, this would call an auth service to validate the token
	// For now, just check if it's our development token
	devToken := os.Getenv("DEV_AUTH_TOKEN")
	if devToken != "" && token == devToken {
		// Return a mock user for development
		return true, &User{
			ID:        "dev-user-1",
			Email:     "dev@example.com",
			Name:      "Development User",
			ProjectID: "dev-project-1",
		}, nil
	}

	// Simple token format for testing: "user_id:email:name:project_id:timestamp"
	// Only for development and testing, not for production
	parts := strings.Split(token, ":")
	if len(parts) == 5 {
		// Check if token is expired (timestamp in seconds)
		if timestamp, err := time.Parse(time.RFC3339, parts[4]); err == nil {
			if time.Now().Before(timestamp) {
				return true, &User{
					ID:        parts[0],
					Email:     parts[1],
					Name:      parts[2],
					ProjectID: parts[3],
				}, nil
			}
		}
	}

	return false, nil, nil
}
