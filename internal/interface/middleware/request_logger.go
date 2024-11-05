package middleware

import (
	"net/http"
	"time"

	"GO-API/internal/pkg/logger"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		logger.Info("Request started: %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)
		logger.Info("Request completed: %s %s, duration=%v", r.Method, r.URL, duration)
	})
}
