package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Processor Service is healthy"))
	})

	log.Println("Starting Processor Service on port 8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
