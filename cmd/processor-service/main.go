package main

import (
	"context"
	"devlog-ai/internal/api/middleware"
	"devlog-ai/internal/messaging"
	"devlog-ai/internal/processor"
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
		port = "8082"
	}

	// Initialize Kafka consumer
	kafkaConsumer, err := messaging.NewKafkaConsumer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	// Create and start the processor service
	processorService := processor.NewService(kafkaConsumer)
	if err := processorService.Start(); err != nil {
		log.Fatalf("Failed to start processor service: %v", err)
	}
	defer processorService.Stop()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", processorService.HealthHandler)

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
		log.Printf("Starting Processor Service on port %s", port)
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
