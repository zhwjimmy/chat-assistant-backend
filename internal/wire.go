//go:build wireinject
// +build wireinject

package internal

import (
	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/handlers"
	"chat-assistant-backend/internal/repositories"
	"chat-assistant-backend/internal/server"
	"chat-assistant-backend/internal/services"

	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.GetDSN()
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// InitializeApp initializes the application with all dependencies
func InitializeApp() (*server.Server, error) {
	wire.Build(
		// Config
		config.Load,

		// Database
		NewDatabase,

		// Repositories
		repositories.NewUserRepository,

		// Services
		services.NewUserService,

		// Handlers
		handlers.NewUserHandler,

		// Server with dependencies
		NewServerWithDependencies,
	)
	return nil, nil
}

// NewServerWithDependencies creates a server with all dependencies injected
func NewServerWithDependencies(
	cfg *config.Config,
	db *gorm.DB,
	userRepo *repositories.UserRepository,
	userService *services.UserService,
	userHandler *handlers.UserHandler,
) *server.Server {
	return server.NewWithDependencies(cfg, db, userHandler)
}
