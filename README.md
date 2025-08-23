# DevLog AI - Backend Services

This directory contains the backend services for the DevLog AI platform.

## Services

- **API Service**: Main API entry point that handles queries and log ingestion
- **Ingestor Service**: Processes incoming logs
- **Processor Service**: Performs analysis on logs
- **ESG Calculator**: Calculates environmental impact metrics

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using the Makefile)

### Building

```bash
# Build all services
make build

# Or build individual services
go build -o bin/api-service ./cmd/api-service
```

### Running

```bash
# Run the API service
make run-api

# Or run directly
go run ./cmd/api-service
```

### Testing

```bash
# Run tests
make test
```

## API Endpoints

### Health Check

```
GET /health
```

### Query Logs

```
POST /api/v1/query

{
  "query": "Show me error logs from the last hour",
  "project_id": "project-123"
}
```

### Ingest Logs

```
POST /api/v1/logs

{
  "project_id": "project-123",
  "logs": [
    {
      "timestamp": "2023-06-01T12:00:00Z",
      "level": "ERROR",
      "message": "Connection failed",
      "source": "app-server-1",
      "service_name": "payment-service"
    }
  ]
}
```
