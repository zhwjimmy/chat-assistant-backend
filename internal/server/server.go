package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/handlers"
	"chat-assistant-backend/internal/logger"
	"chat-assistant-backend/internal/middleware"

	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
	logger *zap.Logger
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

// New creates a new server instance with pre-initialized dependencies
func New(cfg *config.Config, db *gorm.DB, userHandler *handlers.UserHandler, conversationHandler *handlers.ConversationHandler, messageHandler *handlers.MessageHandler, tagHandler *handlers.TagHandler, searchHandler *handlers.SearchHandler) *Server {
	// Set Gin mode
	if cfg.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Add middlewares
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.CORS))

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC(),
			"service":   "chat-assistant-backend",
		})
	})

	// Add Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add API routes
	api := router.Group("/api/v1")
	{
		// User routes
		api.GET("/users/:id", userHandler.GetUser)

		// Tag routes
		api.GET("/tags", tagHandler.GetTags)
		api.GET("/tags/:id", tagHandler.GetTag)
		api.POST("/tags", tagHandler.CreateTag)
		api.PUT("/tags/:id", tagHandler.UpdateTag)
		api.DELETE("/tags/:id", tagHandler.DeleteTag)

		// Conversation routes
		api.GET("/conversations", conversationHandler.GetConversations)
		api.POST("/conversations", conversationHandler.CreateConversation)
		api.GET("/conversations/:id", conversationHandler.GetConversation)
		api.PUT("/conversations/:id/tags", conversationHandler.UpdateConversationTags)
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
