package middleware

import (
	"context"
	"encoding/json"
	"gocleanarchitecture/entities"
	"gocleanarchitecture/interfaces"
	"net/http"
)

// AdminMiddlewareFunc creates a middleware that ensures the user has admin role
func AdminMiddlewareFunc(userRepo interfaces.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// User ID should already be in context from AuthMiddleware
			userID, ok := r.Context().Value("userID").(string)
			if !ok || userID == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Authentication required",
				})
				return
			}

			// Fetch user from repository to check role
			user, err := userRepo.FindByID(userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Failed to verify user permissions",
				})
				return
			}

			if user == nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "User not found",
				})
				return
			}

			// Check if user has admin role
			if user.Role != entities.RoleAdmin {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Admin privileges required",
				})
				return
			}

			// Add user role to context for controllers to use
			ctx := context.WithValue(r.Context(), "userRole", user.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
