package middleware

import (
	"gocleanarchitecture/frameworks/logger"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func LoggingMiddleware(log logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom ResponseWriter to capture the status code
			rw := &responseWriter{w, http.StatusOK}

			// Call the next handler
			next.ServeHTTP(rw, r)

			// Log the request details
			log.Info("HTTP Request",
				logger.LogField{Key: "method", Value: r.Method},
				logger.LogField{Key: "path", Value: r.URL.Path},
				logger.LogField{Key: "status", Value: rw.status},
				logger.LogField{Key: "duration", Value: time.Since(start).String()},
				logger.LogField{Key: "ip", Value: r.RemoteAddr},
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
