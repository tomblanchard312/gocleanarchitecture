package main

import (
	"log"
	"net/http"

	"gocleanarchitecture/config"
	"gocleanarchitecture/frameworks/db"
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web"
	"gocleanarchitecture/interfaces"
	"gocleanarchitecture/usecases"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	customLogger, err := logger.NewLogger(cfg.LogLevel, cfg.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	blogPostRepo := db.NewInMemoryBlogPostRepository()
	blogPostUseCase := &usecases.BlogPostUseCase{Repo: blogPostRepo, Logger: customLogger}
	blogPostController := &interfaces.BlogPostController{BlogPostUseCase: blogPostUseCase}
	router := web.NewRouter(blogPostController, customLogger)

	customLogger.Info("Starting server",
		logger.LogField{Key: "port", Value: cfg.ServerPort},
	)
	err = http.ListenAndServe(cfg.ServerPort, router)
	if err != nil {
		customLogger.Error("Server failed to start",
			logger.LogField{Key: "error", Value: err},
		)
	}
}
