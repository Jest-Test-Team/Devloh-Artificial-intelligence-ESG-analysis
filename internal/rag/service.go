package rag

import (
	"devlog-ai/pkg/apimodels"
	"fmt"
	"log"
	"time"
)

// Service represents the RAG service
type Service struct {
	cloudflareClient *CloudflareClient
}

// NewService creates a new RAG service
func NewService() (*Service, error) {
	cloudflareClient, err := NewCloudflareClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudflare client: %w", err)
	}

	return &Service{
		cloudflareClient: cloudflareClient,
	}, nil
}

// QueryLogs queries logs using RAG
func (s *Service) QueryLogs(request apimodels.QueryRequest) (*apimodels.QueryResponse, error) {
	// Set default options
	options := map[string]interface{}{
		"temperature":      0.0,
		"max_tokens":       1000,
		"include_sources":  true,
		"include_metadata": true,
	}

	// Add any custom options from the request
	if request.Options != nil {
		for k, v := range request.Options {
			options[k] = v
		}
	}

	// Send the query to Cloudflare RAG
	startTime := time.Now()
	response, err := s.cloudflareClient.QueryRAG(request.Query, request.ProjectID, options)
	if err != nil {
		return nil, fmt.Errorf("error querying Cloudflare RAG: %w", err)
	}

	// Extract related logs from the response
	relatedLogs := make([]string, 0, len(response.Citations))
	sources := make([]string, 0, len(response.Citations))
	
	for _, citation := range response.Citations {
		if citation.Text != "" {
			relatedLogs = append(relatedLogs, citation.Text)
		}
		
		// Extract source information from metadata if available
		if citation.Metadata != nil {
			if source, ok := citation.Metadata["source"].(string); ok && source != "" {
				sources = append(sources, source)
			} else if serviceName, ok := citation.Metadata["service_name"].(string); ok && serviceName != "" {
				sources = append(sources, serviceName)
			}
		}
	}

	// Calculate processing time
	processingMs := response.ProcessingMs
	if processingMs == 0 {
		processingMs = int(time.Since(startTime).Milliseconds())
	}

	// Create the response
	result := &apimodels.QueryResponse{
		Answer:       response.Answer,
		RelatedLogs:  relatedLogs,
		Sources:      sources,
		ProcessingMs: int64(processingMs),
	}

	log.Printf("RAG query processed in %d ms: %s", processingMs, request.Query)
	return result, nil
}

// QueryLogsWithContext queries logs using RAG with additional context
func (s *Service) QueryLogsWithContext(request apimodels.QueryRequest, context map[string]interface{}) (*apimodels.QueryResponse, error) {
	// Set default options
	options := map[string]interface{}{
		"temperature":      0.0,
		"max_tokens":       1000,
		"include_sources":  true,
		"include_metadata": true,
	}

	// Add any custom options from the request
	if request.Options != nil {
		for k, v := range request.Options {
			options[k] = v
		}
	}

	// Send the query to Cloudflare RAG with context
	startTime := time.Now()
	response, err := s.cloudflareClient.QueryRAGWithContext(request.Query, request.ProjectID, context, options)
	if err != nil {
		return nil, fmt.Errorf("error querying Cloudflare RAG: %w", err)
	}

	// Extract related logs from the response
	relatedLogs := make([]string, 0, len(response.Citations))
	sources := make([]string, 0, len(response.Citations))
	
	for _, citation := range response.Citations {
		if citation.Text != "" {
			relatedLogs = append(relatedLogs, citation.Text)
		}
		
		// Extract source information from metadata if available
		if citation.Metadata != nil {
			if source, ok := citation.Metadata["source"].(string); ok && source != "" {
				sources = append(sources, source)
			} else if serviceName, ok := citation.Metadata["service_name"].(string); ok && serviceName != "" {
				sources = append(sources, serviceName)
			}
		}
	}

	// Calculate processing time
	processingMs := response.ProcessingMs
	if processingMs == 0 {
		processingMs = int(time.Since(startTime).Milliseconds())
	}

	// Create the response
	result := &apimodels.QueryResponse{
		Answer:       response.Answer,
		RelatedLogs:  relatedLogs,
		Sources:      sources,
		ProcessingMs: int64(processingMs),
	}

	log.Printf("RAG query with context processed in %d ms: %s", processingMs, request.Query)
	return result, nil
}
