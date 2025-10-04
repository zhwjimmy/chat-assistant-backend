package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/handlers"
	"chat-assistant-backend/internal/logger"
	"chat-assistant-backend/internal/repositories"
	"chat-assistant-backend/internal/services"

	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
	logger *zap.Logger
}

// New creates a new server instance
func New(cfg *config.Config, db *gorm.DB) *Server {
	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(loggingMiddleware())
	router.Use(corsMiddleware(cfg.CORS))

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	conversationRepo := repositories.NewConversationRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	conversationService := services.NewConversationService(conversationRepo)
	messageService := services.NewMessageService(messageRepo)
	searchRepo := repositories.NewSearchRepository(db)
	searchService := services.NewSearchService(searchRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	conversationHandler := handlers.NewConversationHandler(conversationService)
	messageHandler := handlers.NewMessageHandler(messageService)
	searchHandler := handlers.NewSearchHandler(searchService)

	// Add health check endpoint
	router.GET("/health", healthCheckHandler)

	// Add API routes
	api := router.Group("/api/v1")
	{
		// User routes
		api.GET("/users/:id", userHandler.GetUser)

		// Conversation routes
		api.GET("/conversations", conversationHandler.GetConversations)
		api.GET("/conversations/:id", conversationHandler.GetConversation)
		api.DELETE("/conversations/:id", conversationHandler.DeleteConversation)
		api.GET("/conversations/:id/messages", messageHandler.GetConversationMessages)

		// Message routes
		api.GET("/messages", messageHandler.GetMessages)
		api.GET("/messages/:id", messageHandler.GetMessage)
		api.DELETE("/messages/:id", messageHandler.DeleteMessage)

		// Search routes
		api.GET("/search", searchHandler.Search)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		config: cfg,
		router: router,
		server: server,
		logger: logger.GetLogger(),
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting HTTP server",
		zap.String("addr", s.server.Addr),
		zap.String("mode", gin.Mode()),
	)

	return s.server.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server...")

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("Failed to shutdown server", zap.Error(err))
		return err
	}

	s.logger.Info("HTTP server stopped")
	return nil
}

// GetRouter returns the Gin router for adding routes
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// requestIDMiddleware adds request ID to each request
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log request
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		requestID, _ := c.Get("request_id")

		logger.WithRequestID(requestID.(string)).Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.Int("body_size", bodySize),
		)
	}
}

// corsMiddleware configures CORS
func corsMiddleware(cfg config.CORSConfig) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}

	return cors.New(corsConfig)
}

// NewWithDependencies creates a new server instance with pre-initialized dependencies
func NewWithDependencies(cfg *config.Config, db *gorm.DB, userHandler *handlers.UserHandler, conversationHandler *handlers.ConversationHandler, messageHandler *handlers.MessageHandler, searchHandler *handlers.SearchHandler) *Server {
	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(requestIDMiddleware())
	router.Use(loggingMiddleware())
	router.Use(corsMiddleware(cfg.CORS))

	// Add health check endpoint
	router.GET("/health", healthCheckHandler)

	// Add Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add API routes
	api := router.Group("/api/v1")
	{
		// User routes
		api.GET("/users/:id", userHandler.GetUser)

		// Conversation routes
		api.GET("/conversations", conversationHandler.GetConversations)
		api.GET("/conversations/:id", conversationHandler.GetConversation)
		api.DELETE("/conversations/:id", conversationHandler.DeleteConversation)
		api.GET("/conversations/:id/messages", messageHandler.GetConversationMessages)

		// Message routes
		api.GET("/messages", messageHandler.GetMessages)
		api.GET("/messages/:id", messageHandler.GetMessage)
		api.DELETE("/messages/:id", messageHandler.DeleteMessage)

		// Search routes
		api.GET("/search", searchHandler.Search)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return &Server{
		config: cfg,
		router: router,
		server: server,
		logger: logger.GetLogger(),
	}
}

// healthCheckHandler handles health check requests
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "chat-assistant-backend",
	})
}
