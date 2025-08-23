package main

import (
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

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	// Create a server with middleware
	server := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
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