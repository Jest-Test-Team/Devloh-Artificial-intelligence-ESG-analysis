package main

import (
	"context"
	"devlog-ai/internal/api/middleware"
	"devlog-ai/internal/ingestor"
	"devlog-ai/internal/messaging"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Initialize Kafka producer
	kafkaProducer, err := messaging.NewKafkaProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Create and start the ingestor service
	ingestorService := ingestor.NewService(kafkaProducer)
	if err := ingestorService.Start(); err != nil {
		log.Fatalf("Failed to start ingestor service: %v", err)
	}
	defer ingestorService.Stop()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})
	mux.HandleFunc("/api/v1/logs", ingestorService.LogHandler)

	// Apply middleware
	handler := middleware.LoggingMiddleware(middleware.CORSMiddleware(mux))

	// Create the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Channel to listen for shutdown signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting Ingestor Service on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-done
	log.Print("Server is shutting down...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Print("Server gracefully stopped")
}