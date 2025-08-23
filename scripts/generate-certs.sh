#!/bin/bash

# Create a directory for certs if it doesn't exist
mkdir -p certs

# Generate a private key
openssl genrsa -out certs/server.key 2048

# Generate a CSR (Certificate Signing Request)
openssl req -new -key certs/server.key -out certs/server.csr -subj "/C=US/ST=State/L=City/O=DevLogAI/OU=Development/CN=localhost"

# Generate a self-signed certificate
openssl x509 -req -days 365 -in certs/server.csr -signkey certs/server.key -out certs/server.crt

echo "Self-signed certificates generated in certs directory"
