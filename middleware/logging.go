package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests with method, path, status, and duration
func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			logger.Info("request handled",
				slog.String("method", r.Method),
				slog.String("path", r.RequestURI),
				slog.Int("status", wrapped.statusCode),
				slog.Duration("duration", duration),
				slog.String("ip", r.RemoteAddr),
			)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures write errors and status code
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == http.StatusOK {
		rw.statusCode = 200
	}
	return rw.ResponseWriter.Write(b)
}
