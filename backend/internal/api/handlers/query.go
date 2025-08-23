package handlers

import (
	"encoding/json"
	"net/http"

	"devlog-ai/internal/rag"
	"devlog-ai/pkg/apimodels"
)

// QueryHandler 處理查詢請求
type QueryHandler struct {
	ragService rag.Service
}

// NewQueryHandler 創建一個新的查詢處理程序
func NewQueryHandler(ragService rag.Service) *QueryHandler {
	return &QueryHandler{
		ragService: ragService,
	}
}

// HandleQuery 處理查詢請求
func (h *QueryHandler) HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var queryRequest apimodels.QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&queryRequest); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 使用 RAG 服務處理查詢
	result, err := h.ragService.Query(queryRequest.Query, queryRequest.Options)
	if err != nil {
		http.Error(w, "Failed to process query: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 構建響應
	response := apimodels.QueryResponse{
		Answer: result,
		ProcessingMs: 0, // 在實際應用中，這裡應該計算處理時間
	}

	// 返回 JSON 響應
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
