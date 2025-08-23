package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Service is healthy"))
	})

	log.Println("Starting API Service on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
