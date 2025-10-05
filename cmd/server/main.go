package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/docs"
	"chat-assistant-backend/internal/logger"
)

// @title Chat Assistant Backend API
// @version 1.0.0
// @description A RESTful API for the Chat Assistant Backend service. This API provides endpoints for managing users, conversations, and messages.
// @contact.name Chat Assistant Team
// @contact.email support@chatassistant.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
func main() {
	// Initialize Swagger docs
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize logger
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Format, cfg.Logging.Output); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	log := logger.GetLogger()
	log.Info("Starting chat-assistant-backend",
		zap.String("version", "1.0.0"),
		zap.String("go_version", "1.23.1"),
	)

	// Initialize application with Wire dependency injection
	app, err := InitializeApp()
	if err != nil {
		log.Fatal("Failed to initialize application", zap.Error(err))
	}

	// Start application
	go func() {
		if err := app.Start(); err != nil {
			log.Error("Failed to start application", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down application...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Shutdown.Timeout)
	defer cancel()

	// Attempt graceful shutdown
	if err := app.Stop(ctx); err != nil {
		log.Error("Application forced to shutdown", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Application exited")
}
