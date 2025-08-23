package middleware

import (
	"net/http"
)

// CORSMiddleware 處理跨域資源共享（CORS）
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 允許的來源
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 允許的方法
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// 允許的標頭
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// 處理預檢請求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// 繼續處理請求
		next.ServeHTTP(w, r)
	})
}