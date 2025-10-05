//go:build wireinject
// +build wireinject

package main

import (
	"chat-assistant-backend/internal/app"
	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/handlers"
	"chat-assistant-backend/internal/infra/database"
	"chat-assistant-backend/internal/infra/elasticsearch"
	"chat-assistant-backend/internal/repositories"
	"chat-assistant-backend/internal/services"

	"github.com/google/wire"
)

// InitializeApp initializes the application with all dependencies
func InitializeApp() (*app.App, error) {
	wire.Build(
		// Config
		config.Load,

		// Infrastructure
		database.DatabaseSet,
		elasticsearch.ElasticsearchSet,

		// Repositories
		repositories.RepositorySet,

		// Services
		services.ServiceSet,

		// Handlers
		handlers.HandlerSet,

		// App with dependencies (includes server.New internally)
		app.New,
	)
	return nil, nil
}
