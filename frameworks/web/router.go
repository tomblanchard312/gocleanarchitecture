package web

import (
	"gocleanarchitecture/frameworks/auth"
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web/middleware"
	"gocleanarchitecture/interfaces"

	"github.com/gorilla/mux"
)

type RouterConfig struct {
	BlogPostController *interfaces.BlogPostController
	AuthController     *interfaces.AuthController
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

	// Blog post routes (public for reading, protected for writing)
	router.HandleFunc("/blogposts", config.BlogPostController.GetAllBlogPosts).Methods("GET")
	router.HandleFunc("/blogposts/{id}", config.BlogPostController.GetBlogPost).Methods("GET")

	// Protected blog post routes
	protectedBlogRouter := router.PathPrefix("/blogposts").Subrouter()
	protectedBlogRouter.Use(middleware.AuthMiddlewareFunc(config.JWTManager))
	protectedBlogRouter.HandleFunc("", config.BlogPostController.CreateBlogPost).Methods("POST")
	protectedBlogRouter.HandleFunc("/{id}", config.BlogPostController.UpdateBlogPost).Methods("PUT")
	protectedBlogRouter.HandleFunc("/{id}", config.BlogPostController.DeleteBlogPost).Methods("DELETE")

	return router
}
