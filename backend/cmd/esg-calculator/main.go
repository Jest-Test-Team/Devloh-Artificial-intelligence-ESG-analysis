package main

import (
	"log"
	"net/http"
)

func main() {
	// Simple health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ESG Calculator Service is healthy"))
	})

	log.Println("Starting ESG Calculator Service on port 8083")
	if err := http.ListenAndServe(":8083", nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
