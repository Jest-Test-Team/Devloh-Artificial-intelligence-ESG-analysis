package main

import (
	"context"
	"crypto/tls"
	"devlog-ai/internal/api/middleware"
	"devlog-ai/internal/embedding"
	"devlog-ai/internal/messaging"
	"devlog-ai/internal/processor"
	"devlog-ai/internal/vectordb"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Check if vector database integration should be enabled
	enableVectorDB := true
	if val := os.Getenv("ENABLE_VECTOR_DB"); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			enableVectorDB = b
		}
	}

	// Initialize Kafka consumer
	kafkaConsumer, err := messaging.NewKafkaConsumer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	var embeddingClient embedding.Client
	var vectorDBClient vectordb.Client

	// Initialize embedding client if vector database is enabled
	if enableVectorDB {
		// Get embedding client configuration
		embeddingConfig, err := embedding.GetConfigFromEnv()
		if err != nil {
			log.Printf("Warning: %v", err)
			log.Printf("Embedding client will not be available")
		} else {
			embeddingClient, err = embedding.NewClient(embeddingConfig)
			if err != nil {
				log.Printf("Warning: Failed to initialize embedding client: %v", err)
				log.Printf("Embedding client will not be available")
			} else {
				log.Printf("Embedding client initialized with type: %s", embeddingConfig.Type)
			}
		}

		// Get vector database configuration
		vectorDBConfig, err := vectordb.GetConfigFromEnv()
		if err != nil {
			log.Printf("Warning: %v", err)
			log.Printf("Vector database will not be available")
		} else {
			vectorDBClient, err = vectordb.NewClient(vectorDBConfig)
			if err != nil {
				log.Printf("Warning: Failed to initialize vector database client: %v", err)
				log.Printf("Vector database will not be available")
			} else {
				log.Printf("Vector database client initialized with type: %s", vectorDBConfig.Type)
				
				// Test connection to vector database
				if err := vectorDBClient.Connect(); err != nil {
					log.Printf("Warning: Failed to connect to vector database: %v", err)
					log.Printf("Vector database will not be available")
					vectorDBClient = nil
				} else {
					log.Printf("Successfully connected to vector database")
				}
			}
		}
	}

	// Create and start the processor service
	processorService := processor.NewService(
		kafkaConsumer,
		embeddingClient,
		vectorDBClient,
		processor.ServiceConfig{
			EnableVectorDB: enableVectorDB && embeddingClient != nil && vectorDBClient != nil,
		},
	)
	
	if err := processorService.Start(); err != nil {
		log.Fatalf("Failed to start processor service: %v", err)
	}
	defer processorService.Stop()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define public paths that don't require authentication
	publicPaths := []string{
		"/health",
		// Add other public paths as needed
	}

	// Register routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

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
			log.Printf("Starting Processor Service on HTTP port %s", port)
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("Error starting server: %v", err)
				cancel()
			}
		} else {
			log.Printf("Starting Processor Service on HTTPS port %s", port)
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