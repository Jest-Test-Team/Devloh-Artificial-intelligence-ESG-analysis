package main

import (
	"crypto/tls"
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
		port = "8081"
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

	// Check if TLS certificates exist
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		log.Printf("Warning: TLS certificate not found at %s, falling back to HTTP", certPath)
		log.Printf("Starting Ingestor Service on HTTP port %s", port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	} else {
		log.Printf("Starting Ingestor Service on HTTPS port %s", port)
		if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}
}