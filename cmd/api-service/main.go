package main

import (
	"context"
	"crypto/tls"
	"devlog-ai/internal/api/handlers"
	"devlog-ai/internal/api/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
		log.Printf("Warning: Failed to initialize log handlers: %v", err)
	}

	if err := handlers.InitQueryHandler(); err != nil {
		log.Printf("Warning: Failed to initialize query handler: %v", err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define public paths that don't require authentication
	publicPaths := []string{
		"/health",
		// Add other public paths as needed
	}

	// Register routes
	mux.HandleFunc("/health", handlers.HealthHandler)
	mux.HandleFunc("/api/v1/logs", handlers.IngestLogsHandler)
	mux.HandleFunc("/api/v1/query", handlers.QueryHandler)

	// Apply middleware
	handler := middleware.LoggingMiddleware(middleware.CORSMiddleware(middleware.AuthMiddleware(publicPaths)(mux)))

	// Create a server with TLS
	server := http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}

	certPath := os.Getenv("TLS_CERT_PATH")
	keyPath := os.Getenv("TLS_KEY_PATH")

	// Default to looking in a certs directory if not specified
	if certPath == "" {
		certPath = filepath.Join("certs", "server.crt")
	}
	if keyPath == "" {
		keyPath = filepath.Join("certs", "server.key")
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Check if TLS certificates exist
		if _, err := os.Stat(certPath); os.IsNotExist(err) {
			log.Printf("Warning: TLS certificate not found at %s, falling back to HTTP", certPath)
			log.Printf("Starting API Service on HTTP port %s", port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Error starting server: %v", err)
				cancel()
			}
		} else {
			log.Printf("Starting API Service on HTTPS port %s", port)
			if err := server.ListenAndServeTLS(certPath, keyPath); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Error starting server: %v", err)
				cancel()
			}
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	log.Print("Server shutting down...")

	// Create a deadline for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Print("Server gracefully stopped")
}