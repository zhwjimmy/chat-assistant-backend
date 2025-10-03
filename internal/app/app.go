package app

import (
	"context"
	"time"

	"go.uber.org/zap"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/logger"
	"chat-assistant-backend/internal/server"
)

// App represents the application
type App struct {
	config *config.Config
	server *server.Server
	logger *zap.Logger
}

// New creates a new application instance
func New(cfg *config.Config) *App {
	return &App{
		config: cfg,
		server: server.New(cfg),
		logger: logger.GetLogger(),
	}
}

// Start starts the application
func (a *App) Start() error {
	a.logger.Info("Starting application...")

	// Start server in a goroutine
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Error("Failed to start server", zap.Error(err))
		}
	}()

	return nil
}

// Stop gracefully stops the application
func (a *App) Stop(ctx context.Context) error {
	a.logger.Info("Stopping application...")

	// Stop server
	if err := a.server.Stop(ctx); err != nil {
		a.logger.Error("Failed to stop server", zap.Error(err))
		return err
	}

	a.logger.Info("Application stopped")
	return nil
}

// GetServer returns the server instance
func (a *App) GetServer() *server.Server {
	return a.server
}
