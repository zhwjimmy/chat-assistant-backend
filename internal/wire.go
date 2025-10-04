//go:build wireinject
// +build wireinject

package internal

import (
	"fmt"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/handlers"
	"chat-assistant-backend/internal/infra/elasticsearch"
	"chat-assistant-backend/internal/migrations"
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewElasticsearchClient creates a new Elasticsearch client
func NewElasticsearchClient(cfg *config.Config) (*elasticsearch.Client, error) {
	esConfig := &elasticsearch.Config{
		Hosts:    cfg.Elasticsearch.Hosts,
		Username: cfg.Elasticsearch.Username,
		Password: cfg.Elasticsearch.Password,
		Timeout:  cfg.Elasticsearch.Timeout,
		Index: elasticsearch.IndexConfig{
			Conversations: cfg.Elasticsearch.Index.Conversations,
			Messages:      cfg.Elasticsearch.Index.Messages,
		},
	}

	return elasticsearch.NewClient(esConfig)
}

// NewElasticsearchIndexer creates a new Elasticsearch indexer
func NewElasticsearchIndexer(esClient *elasticsearch.Client, cfg *config.Config) repositories.ElasticsearchIndexer {
	return repositories.NewElasticsearchIndexer(esClient.GetClient(), cfg.Elasticsearch.Index.Conversations)
}

// NewElasticsearchRepository creates a new Elasticsearch repository
func NewElasticsearchRepository(esClient *elasticsearch.Client, cfg *config.Config) services.SearchRepository {
	return repositories.NewElasticsearchRepository(esClient.GetClient(), cfg.Elasticsearch.Index.Conversations)
}

// NewSearchRepository creates a search repository with Elasticsearch
func NewSearchRepository(cfg *config.Config) services.SearchRepository {
	// Create Elasticsearch client
	esClient, err := NewElasticsearchClient(cfg)
	if err != nil {
		// If Elasticsearch is not available, panic - search is required
		panic(fmt.Sprintf("Failed to initialize Elasticsearch client: %v", err))
	}

	// Use Elasticsearch implementation
	return repositories.NewElasticsearchRepository(esClient.GetClient(), cfg.Elasticsearch.Index.Conversations)
}

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	migrator, err := migrations.NewMigrator(db, nil)
	if err != nil {
		return err
	}

	return migrator.Up()
}

// InitializeDatabase creates database connection and runs migrations
func InitializeDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

// InitializeApp initializes the application with all dependencies
func InitializeApp() (*server.Server, error) {
	wire.Build(
		// Config
		config.Load,

		// Infrastructure
		InitializeDatabase,

		// Repositories
		repositories.NewUserRepository,
		repositories.NewConversationRepository,
		repositories.NewMessageRepository,
		NewSearchRepository,

		// Services
		services.NewUserService,
		services.NewConversationService,
		services.NewMessageService,
		services.NewSearchService,

		// Handlers
		handlers.NewUserHandler,
		handlers.NewConversationHandler,
		handlers.NewMessageHandler,
		handlers.NewSearchHandler,

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
	conversationRepo *repositories.ConversationRepository,
	messageRepo *repositories.MessageRepository,
	searchRepo services.SearchRepository,
	userService *services.UserService,
	conversationService *services.ConversationService,
	messageService *services.MessageService,
	searchService *services.SearchService,
	userHandler *handlers.UserHandler,
	conversationHandler *handlers.ConversationHandler,
	messageHandler *handlers.MessageHandler,
	searchHandler *handlers.SearchHandler,
) *server.Server {
	return server.NewWithDependencies(cfg, db, userHandler, conversationHandler, messageHandler, searchHandler)
}
