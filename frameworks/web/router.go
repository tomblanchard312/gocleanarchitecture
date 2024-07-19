package web

import (
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web/middleware"
	"gocleanarchitecture/interfaces"

	"github.com/gorilla/mux"
)

func NewRouter(controller *interfaces.BlogPostController, logger logger.Logger) *mux.Router {
	router := mux.NewRouter()

	// Add middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))

	// Define routes
	router.HandleFunc("/blogposts", controller.CreateBlogPost).Methods("POST")
	router.HandleFunc("/blogposts", controller.GetAllBlogPosts).Methods("GET")
	router.HandleFunc("/blogposts/{id}", controller.GetBlogPost).Methods("GET")
	router.HandleFunc("/blogposts/{id}", controller.UpdateBlogPost).Methods("PUT")
	router.HandleFunc("/blogposts/{id}", controller.DeleteBlogPost).Methods("DELETE")

	return router
}
