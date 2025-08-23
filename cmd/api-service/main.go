package main

import (
	"devlog-ai/internal/api/handlers"
	"devlog-ai/internal/api/middleware"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize handlers
	if err := handlers.InitLogHandlers(); err != nil {
		log.Fatalf("Failed to initialize log handlers: %v", err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/api/v1/query", handlers.QueryHandler)
	mux.HandleFunc("/api/v1/logs", handlers.IngestLogsHandler)

	// Create a server with middleware
	server := http.Server{
		Addr:         ":" + port,
		Handler:      applyMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server
	log.Printf("Starting API Service on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func applyMiddleware(handler http.Handler) http.Handler {
	// Apply middleware in reverse order
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)
	return handler
}