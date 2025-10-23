package main

import (
	"gocleanarchitecture/config"
	"gocleanarchitecture/frameworks/auth"
	"gocleanarchitecture/frameworks/db"
	"gocleanarchitecture/frameworks/db/sqlite"
	"gocleanarchitecture/frameworks/db/supabase"
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web"
	"gocleanarchitecture/interfaces"
	"gocleanarchitecture/usecases"
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	customLogger, err := logger.NewLogger(cfg.LogLevel, cfg.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Create repositories based on configuration
	var blogPostRepo interfaces.BlogPostRepository
	var userRepo interfaces.UserRepository

	switch strings.ToLower(cfg.DBType) {
	case "supabase":
		if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
			log.Fatal("Supabase URL and Key must be provided when using Supabase")
		}
		blogPostRepo = supabase.NewSupabaseBlogPostRepository(cfg.SupabaseURL, cfg.SupabaseKey)
		userRepo = supabase.NewSupabaseUserRepository(cfg.SupabaseURL, cfg.SupabaseKey)
		customLogger.Info("Using Supabase repository", logger.Field("url", cfg.SupabaseURL))
	case "inmemory":
		blogPostRepo = db.NewInMemoryBlogPostRepository()
		customLogger.Info("Using in-memory repository")
		log.Println("Warning: User authentication not supported with in-memory database")
	case "sqlite":
		fallthrough
	default:
		sqliteDB, err := sqlite.InitDB(cfg.DBPath)
		if err != nil {
			log.Fatalf("Failed to initialize SQLite database: %v", err)
		}
		defer sqliteDB.Close()
		blogPostRepo = sqlite.NewSQLiteBlogPostRepository(sqliteDB)
		customLogger.Info("Using SQLite repository", logger.Field("path", cfg.DBPath))
		log.Println("Warning: User authentication not fully supported with SQLite (requires schema migration)")
	}

	// Initialize JWT Manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTTokenDuration)
	tokenGenerator := auth.NewTokenGeneratorAdapter(jwtManager)

	// Create adapters and use cases with proper dependency injection
	useCaseLogger := logger.NewUseCaseLoggerAdapter(customLogger)

	// Blog post use case
	blogPostUseCase := usecases.NewBlogPostUseCase(blogPostRepo, useCaseLogger)
	blogPostController := &interfaces.BlogPostController{BlogPostUseCase: blogPostUseCase}

	// Auth use case (only if user repository is available)
	var authController *interfaces.AuthController
	if userRepo != nil {
		authUseCase := usecases.NewAuthUseCase(userRepo, tokenGenerator, useCaseLogger)
		authController = &interfaces.AuthController{AuthUseCase: authUseCase}
		customLogger.Info("Authentication enabled")
	} else {
		customLogger.Warn("Authentication disabled - user repository not available")
	}

	// Create router with all controllers
	routerConfig := &web.RouterConfig{
		BlogPostController: blogPostController,
		AuthController:     authController,
		JWTManager:         jwtManager,
		Logger:             customLogger,
	}
	router := web.NewRouter(routerConfig)

	customLogger.Info("Starting server", logger.Field("port", cfg.ServerPort))
	log.Fatal(http.ListenAndServe(cfg.ServerPort, router))
}
