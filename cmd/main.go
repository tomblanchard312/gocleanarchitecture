package main

import (
	"gocleanarchitecture/config"
	"gocleanarchitecture/frameworks/db/sqlite"
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web"
	"gocleanarchitecture/interfaces"
	"gocleanarchitecture/usecases"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := sqlite.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	customLogger, err := logger.NewLogger(cfg.LogLevel, cfg.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	repo := sqlite.NewSQLiteBlogPostRepository(db)
	blogPostUseCase := &usecases.BlogPostUseCase{Repo: repo, Logger: customLogger}
	controller := &interfaces.BlogPostController{BlogPostUseCase: blogPostUseCase}

	router := web.NewRouter(controller, customLogger)

	customLogger.Info("Starting server", logger.Field("port", cfg.ServerPort))
	log.Fatal(http.ListenAndServe(cfg.ServerPort, router))
}
