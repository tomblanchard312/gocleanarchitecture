package middleware

import (
	"gocleanarchitecture/frameworks/logger"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

func RecoveryMiddleware(log logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the error and stack trace
					log.Error("Panic recovered",
						logger.LogField{Key: "error", Value: err},
						logger.LogField{Key: "stack", Value: string(debug.Stack())},
					)

					// Return an internal server error to the client
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
