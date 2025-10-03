package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/logger"
)

// Server represents the HTTP server
type Server struct {
	config *config.Config
	router *gin.Engine
	server *http.Server
	logger *zap.Logger
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
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

// healthCheckHandler handles health check requests
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().UTC(),
		"service":   "chat-assistant-backend",
	})
}
