package main

import (
	"gocleanarchitecture/config"
	"gocleanarchitecture/frameworks/auth"
	"gocleanarchitecture/frameworks/db"
	"gocleanarchitecture/frameworks/db/sqlite"
	"gocleanarchitecture/frameworks/db/supabase"
	"gocleanarchitecture/frameworks/logger"
	"gocleanarchitecture/frameworks/web"
	"gocleanarchitecture/frameworks/websocket"
	"gocleanarchitecture/interfaces"
	"gocleanarchitecture/usecases"

	"log"
	"net/http"
	"strings"

	"github.com/rs/cors"

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
	var commentRepo interfaces.CommentRepository

	switch strings.ToLower(cfg.DBType) {
	case "supabase":
		if cfg.SupabaseURL == "" || cfg.SupabaseKey == "" {
			log.Fatal("Supabase URL and Key must be provided when using Supabase")
		}
		blogPostRepo = supabase.NewSupabaseBlogPostRepository(cfg.SupabaseURL, cfg.SupabaseKey)
		userRepo = supabase.NewSupabaseUserRepository(cfg.SupabaseURL, cfg.SupabaseKey)
		commentRepo = supabase.NewSupabaseCommentRepository(cfg.SupabaseURL, cfg.SupabaseKey)
		customLogger.Info("Using Supabase repository", logger.Field("url", cfg.SupabaseURL))
	case "inmemory":
		blogPostRepo = db.NewInMemoryBlogPostRepository()
		userRepo = db.NewInMemoryUserRepository()
		commentRepo = db.NewInMemoryCommentRepository()
		customLogger.Info("Using in-memory repository")
		customLogger.Warn("In-memory database: data will be lost on restart")
	case "sqlite":
		fallthrough
	default:
		sqliteDB, err := sqlite.InitDB(cfg.DBPath)
		if err != nil {
			log.Fatalf("Failed to initialize SQLite database: %v", err)
		}
		defer sqliteDB.Close()
		blogPostRepo = sqlite.NewSQLiteBlogPostRepository(sqliteDB)
		userRepo = sqlite.NewSQLiteUserRepository(sqliteDB)
		commentRepo = sqlite.NewSQLiteCommentRepository(sqliteDB)
		customLogger.Info("Using SQLite repository", logger.Field("path", cfg.DBPath))
	}

	// Initialize JWT Manager
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTTokenDuration)
	tokenGenerator := auth.NewTokenGeneratorAdapter(jwtManager)

	// Create adapters and use cases with proper dependency injection
	useCaseLogger := logger.NewUseCaseLoggerAdapter(customLogger)

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run() // Start hub in a goroutine
	customLogger.Info("WebSocket hub started")

	// Blog post use case
	blogPostUseCase := usecases.NewBlogPostUseCase(blogPostRepo, useCaseLogger)
	blogPostController := &interfaces.BlogPostController{
		BlogPostUseCase: blogPostUseCase,
		WebSocketHub:    wsHub,
	}

	// Comment use case
	commentUseCase := usecases.NewCommentUseCase(commentRepo, blogPostRepo, userRepo, useCaseLogger)
	commentController := &interfaces.CommentController{
		CommentUseCase: commentUseCase,
		WebSocketHub:   wsHub,
	}

	// WebSocket handler
	wsHandler := interfaces.NewWebSocketHandler(wsHub)

	// Auth use case (only if user repository is available)
	var authController *interfaces.AuthController
	var adminController *interfaces.AdminController
	var oauth2Controller *interfaces.OAuth2Controller
	if userRepo != nil {
		authUseCase := usecases.NewAuthUseCase(userRepo, tokenGenerator, useCaseLogger)
		authController = &interfaces.AuthController{AuthUseCase: authUseCase}

		// Admin use case
		adminUseCase := usecases.NewAdminUseCase(userRepo, useCaseLogger)
		adminController = interfaces.NewAdminController(adminUseCase)

		// OAuth2 providers (optional - only if configured)
		var googleProvider *auth.OAuth2Provider
		var githubProvider *auth.OAuth2Provider

		if cfg.GoogleClientID != "" && cfg.GoogleClientSecret != "" {
			googleProvider = auth.NewGoogleOAuth2Provider(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL)
			customLogger.Info("Google OAuth2 enabled")
		}

		if cfg.GitHubClientID != "" && cfg.GitHubClientSecret != "" {
			githubProvider = auth.NewGitHubOAuth2Provider(cfg.GitHubClientID, cfg.GitHubClientSecret, cfg.GitHubRedirectURL)
			customLogger.Info("GitHub OAuth2 enabled")
		}

		// Create OAuth2 controller if at least one provider is configured
		if googleProvider != nil || githubProvider != nil {
			oauth2Controller = interfaces.NewOAuth2Controller(googleProvider, githubProvider, authUseCase)
			customLogger.Info("OAuth2 social login enabled")
		}

		customLogger.Info("Authentication and admin features enabled")
	} else {
		customLogger.Warn("Authentication disabled - user repository not available")
	}

	// Create router with all controllers
	routerConfig := &web.RouterConfig{
		BlogPostController: blogPostController,
		AuthController:     authController,
		AdminController:    adminController,
		CommentController:  commentController,
		WebSocketHandler:   wsHandler,
		OAuth2Controller:   oauth2Controller,
		UserRepo:           userRepo,
		JWTManager:         jwtManager,
		Logger:             customLogger,
	}
	router := web.NewRouter(routerConfig)

	// Add CORS middleware with WebSocket support
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		// Allow WebSocket upgrade headers
		ExposedHeaders: []string{"*"},
	})

	handler := corsHandler.Handler(router)

	customLogger.Info("Starting server", logger.Field("port", cfg.ServerPort))
	log.Fatal(http.ListenAndServe(cfg.ServerPort, handler))
}
