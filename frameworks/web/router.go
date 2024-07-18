package web

import (
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web/middleware"
	"gocleanarchitecture/interfaces"

	"github.com/gorilla/mux"
)

func NewRouter(blogPostController *interfaces.BlogPostController, log logger.Logger) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware(log))
	router.Use(middleware.RecoveryMiddleware(log))

	router.HandleFunc("/blogposts", blogPostController.CreateBlogPost).Methods("POST")
	router.HandleFunc("/blogposts", blogPostController.GetAllBlogPosts).Methods("GET")

	return router
}
