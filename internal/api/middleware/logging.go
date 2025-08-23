package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs information about each request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Process the request
		next.ServeHTTP(w, r)

		// Log request details
		log.Printf(
			"Method: %s | Path: %s | Duration: %v | User-Agent: %s",
			r.Method,
			r.URL.Path,
			time.Since(start),
			r.UserAgent(),
		)
	})
}