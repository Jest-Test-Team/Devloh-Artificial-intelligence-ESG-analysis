.PHONY: build run test clean fmt vet

# Default target
all: test build

# Build all services
build:
	@echo "Building API service..."
	go build -o bin/api-service ./cmd/api-service
	@echo "Building Ingestor service..."
	go build -o bin/ingestor-service ./cmd/ingestor-service
	@echo "Building Processor service..."
	go build -o bin/processor-service ./cmd/processor-service
	@echo "Building ESG Calculator service..."
	go build -o bin/esg-calculator ./cmd/esg-calculator

# Build and run the API service
run-api:
	@echo "Running API service..."
	go run ./cmd/api-service

# Build and run the Ingestor service
run-ingestor:
	@echo "Running Ingestor service..."
	go run ./cmd/ingestor-service

# Build and run the Processor service
run-processor:
	@echo "Running Processor service..."
	go run ./cmd/processor-service

# Build and run the ESG Calculator service
run-esg:
	@echo "Running ESG Calculator service..."
	go run ./cmd/esg-calculator

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run static analysis
vet:
	@echo "Vetting code..."
	go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Done."
