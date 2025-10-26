package web

import (
	"gocleanarchitecture/frameworks/auth"
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web/middleware"
	"gocleanarchitecture/interfaces"
	"net/http"

	"github.com/gorilla/mux"
)

type RouterConfig struct {
	BlogPostController *interfaces.BlogPostController
	AuthController     *interfaces.AuthController
	AdminController    *interfaces.AdminController
	CommentController  *interfaces.CommentController
	WebSocketHandler   *interfaces.WebSocketHandler
	OAuth2Controller   *interfaces.OAuth2Controller
	UserRepo           interfaces.UserRepository
	JWTManager         *auth.JWTManager
	Logger             logger.Logger
}

func NewRouter(config *RouterConfig) *mux.Router {
	router := mux.NewRouter()

	// Add global middleware
	router.Use(middleware.LoggingMiddleware(config.Logger))
	router.Use(middleware.RecoveryMiddleware(config.Logger))

	// Auth routes (public - no authentication required)
	authRouter := router.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", config.AuthController.Register).Methods("POST")
	authRouter.HandleFunc("/login", config.AuthController.Login).Methods("POST")
	authRouter.HandleFunc("/users/{username}", config.AuthController.GetUserByUsername).Methods("GET")

	// Protected auth routes (requires authentication)
	protectedAuthRouter := router.PathPrefix("/auth").Subrouter()
	protectedAuthRouter.Use(middleware.AuthMiddlewareFunc(config.JWTManager))
	protectedAuthRouter.HandleFunc("/profile", config.AuthController.GetProfile).Methods("GET")
	protectedAuthRouter.HandleFunc("/profile", config.AuthController.UpdateProfile).Methods("PUT")
	protectedAuthRouter.HandleFunc("/change-password", config.AuthController.ChangePassword).Methods("POST")

	// OAuth2 routes (public - no authentication required)
	if config.OAuth2Controller != nil {
		authRouter.HandleFunc("/google", config.OAuth2Controller.InitiateGoogleLogin).Methods("GET")
		authRouter.HandleFunc("/google/callback", config.OAuth2Controller.GoogleCallback).Methods("GET")
		authRouter.HandleFunc("/github", config.OAuth2Controller.InitiateGitHubLogin).Methods("GET")
		authRouter.HandleFunc("/github/callback", config.OAuth2Controller.GitHubCallback).Methods("GET")
	}

	// Blog post routes (public for reading, protected for writing)
	router.HandleFunc("/blogposts", config.BlogPostController.GetAllBlogPosts).Methods("GET")
	router.HandleFunc("/blogposts/{id}", config.BlogPostController.GetBlogPost).Methods("GET")

	// Protected blog post routes
	protectedBlogRouter := router.PathPrefix("/blogposts").Subrouter()
	protectedBlogRouter.Use(middleware.AuthMiddlewareFunc(config.JWTManager))
	protectedBlogRouter.HandleFunc("", config.BlogPostController.CreateBlogPost).Methods("POST")
	protectedBlogRouter.HandleFunc("/{id}", config.BlogPostController.UpdateBlogPost).Methods("PUT")
	protectedBlogRouter.HandleFunc("/{id}", config.BlogPostController.DeleteBlogPost).Methods("DELETE")

	// Admin routes (requires authentication + admin role)
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddlewareFunc(config.JWTManager))
	adminRouter.Use(middleware.AdminMiddlewareFunc(config.UserRepo))
	adminRouter.HandleFunc("/users", config.AdminController.GetAllUsers).Methods("GET")
	adminRouter.HandleFunc("/users/{id}", config.AdminController.GetUserDetails).Methods("GET")
	adminRouter.HandleFunc("/users/{id}/role", config.AdminController.UpdateUserRole).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}", config.AdminController.DeleteUser).Methods("DELETE")

	// Comment routes (public for reading, protected for writing)
	router.HandleFunc("/blogposts/{blogPostId}/comments", config.CommentController.GetCommentsByBlogPost).Methods("GET")
	router.HandleFunc("/comments/{commentId}/replies", config.CommentController.GetRepliesByComment).Methods("GET")

	// Protected comment routes
	protectedCommentRouter := router.PathPrefix("").Subrouter()
	protectedCommentRouter.Use(middleware.AuthMiddlewareFunc(config.JWTManager))
	protectedCommentRouter.HandleFunc("/blogposts/{blogPostId}/comments", config.CommentController.CreateComment).Methods("POST")
	protectedCommentRouter.HandleFunc("/comments/{commentId}", config.CommentController.UpdateComment).Methods("PUT")
	protectedCommentRouter.HandleFunc("/comments/{commentId}", config.CommentController.DeleteComment).Methods("DELETE")

	// WebSocket endpoint (public - can be accessed by anyone)
	router.HandleFunc("/ws", config.WebSocketHandler.HandleWebSocket).Methods("GET")

	// Swagger/API Documentation endpoint
	router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./api-documentation.yaml")
	}).Methods("GET")

	// Swagger UI redirect
	router.HandleFunc("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://editor.swagger.io/?url=http://localhost:8080/swagger", http.StatusTemporaryRedirect)
	}).Methods("GET")

	return router
}
