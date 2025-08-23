package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// LoggingMiddleware 是一個記錄 HTTP 請求的中間件
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 開始時間
		start := time.Now()
		
		// 創建一個響應記錄器，以便我們可以捕獲狀態碼
		wrapped := wrapResponseWriter(w)
		
		// 調用下一個處理器
		next.ServeHTTP(wrapped, r)
		
		// 記錄請求
		clientIP := getClientIP(r)
		method := r.Method
		path := r.URL.RequestURI()
		statusCode := wrapped.status
		duration := time.Since(start)
		userAgent := r.UserAgent()
		referer := r.Referer()
		
		log.Printf(
			"[%s] %s %s %d %v %s %s",
			method,
			path,
			clientIP,
			statusCode,
			duration,
			userAgent,
			referer,
		)
	})
}

// responseWriter 是 http.ResponseWriter 的包裝器
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

// 包裝 ResponseWriter
func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, status: http.StatusOK}
}

// WriteHeader 覆蓋以捕獲狀態碼
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// Write 覆蓋以確保狀態碼被設置
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// getClientIP 嘗試獲取客戶端的真實 IP 地址
func getClientIP(r *http.Request) string {
	// 檢查代理標頭
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		if addresses := r.Header.Get(h); addresses != "" {
			// X-Forwarded-For 可能包含多個 IP，取第一個作為客戶端 IP
			parts := strings.Split(addresses, ",")
			return strings.TrimSpace(parts[0])
		}
	}
	
	// 沒有代理標頭，嘗試從 RemoteAddr 獲取
	ip := r.RemoteAddr
	
	// 如果 RemoteAddr 包含端口，去掉端口部分
	if i := strings.LastIndex(ip, ":"); i != -1 {
		ip = ip[:i]
	}
	
	return ip
}