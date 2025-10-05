# Chat Assistant Backend Makefile

# Variables
BINARY_NAME=chat-assistant-backend
DOCKER_IMAGE=chat-assistant-backend
DOCKER_TAG=latest
GO_VERSION=1.23.1

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build parameters
BUILD_DIR=./bin
MAIN_PATH=./cmd/server
IMPORTER_BINARY=chat-assistant-importer
IMPORTER_PATH=./cmd/importer
MIGRATE_BINARY=chat-assistant-migrate
MIGRATE_PATH=./cmd/migrate
ES_MANAGER_BINARY=chat-assistant-es-manager
ES_MANAGER_PATH=./cmd/es-manager
DATA_SYNC_BINARY=chat-assistant-data-sync
DATA_SYNC_PATH=./cmd/data-sync

# Migration parameters
MIGRATIONS_DIR=./internal/migrations

.PHONY: all build clean test deps run docker-build docker-run gen-swagger gen-wire lint help dev-db-up dev-db-down dev-db-logs dev-db-reset dev-setup dev-clean build-importer run-importer test-import build-migrate migrate-up migrate-down migrate-reset migrate-status migrate-version migrate-create migrate-fix migrate-validate build-es-manager es-status es-init es-recreate es-health build-data-sync sync-data sync-data-dry

# Default target
all: deps build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Build importer tool
build-importer:
	@echo "Building $(IMPORTER_BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(IMPORTER_BINARY) -v $(IMPORTER_PATH)
	@echo "Importer build completed: $(BUILD_DIR)/$(IMPORTER_BINARY)"

# Build migration tool
build-migrate:
	@echo "Building $(MIGRATE_BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(MIGRATE_BINARY) -v $(MIGRATE_PATH)
	@echo "Migration tool build completed: $(BUILD_DIR)/$(MIGRATE_BINARY)"

# Build ES manager tool
build-es-manager:
	@echo "Building $(ES_MANAGER_BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(ES_MANAGER_BINARY) -v $(ES_MANAGER_PATH)
	@echo "ES Manager build completed: $(BUILD_DIR)/$(ES_MANAGER_BINARY)"

# Build data sync tool
build-data-sync:
	@echo "Building $(DATA_SYNC_BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(DATA_SYNC_BINARY) -v $(DATA_SYNC_PATH)
	@echo "Data Sync tool build completed: $(BUILD_DIR)/$(DATA_SYNC_BINARY)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run importer tool
run-importer:
	@echo "Running importer tool..."
	$(GOCMD) run $(IMPORTER_PATH) $(ARGS)

# Test import functionality
test-import:
	@echo "Running import tests..."
	$(GOTEST) -v ./internal/importer/...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy


# Run the application locally
run:
	@echo "Running $(BINARY_NAME)..."
	$(GOCMD) run $(MAIN_PATH)

# Run with hot reload (requires air)
run-dev:
	@echo "Running with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		$(GOCMD) run $(MAIN_PATH); \
	fi

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker compose up (PostgreSQL only)
docker-compose-up:
	@echo "Starting PostgreSQL with docker-compose..."
	docker-compose up -d

# Docker compose down
docker-compose-down:
	@echo "Stopping PostgreSQL..."
	docker-compose down

# Docker compose with migration
docker-compose-up-migrate:
	@echo "Starting PostgreSQL and running migrations..."
	docker-compose --profile migrate up -d

# Build migration Docker image
docker-build-migrate:
	@echo "Building migration Docker image..."
	docker build -f Dockerfile.migrate -t $(DOCKER_IMAGE)-migrate:$(DOCKER_TAG) .
	@echo "Migration Docker image built: $(DOCKER_IMAGE)-migrate:$(DOCKER_TAG)"

# Run migration in Docker
docker-migrate-up:
	@echo "Running migrations in Docker..."
	docker-compose --profile migrate run --rm migrate migrate -command up

# Run migration rollback in Docker
docker-migrate-down:
	@echo "Rolling back migrations in Docker..."
	docker-compose --profile migrate run --rm migrate migrate -command down

# Development database management
dev-db-up:
	@echo "Starting PostgreSQL for local development..."
	docker-compose up postgres -d
	@echo "PostgreSQL started. Connect with: postgres://postgres:postgres@localhost:5432/chat_assistant"

dev-db-down:
	@echo "Stopping PostgreSQL..."
	docker-compose down

dev-db-logs:
	@echo "Showing PostgreSQL logs..."
	docker-compose logs -f postgres

dev-db-reset:
	@echo "Resetting PostgreSQL database..."
	docker-compose down
	docker volume rm chat-assistant-backend_postgres_data 2>/dev/null || true
	docker-compose up postgres -d
	@echo "Database reset completed"

# Local development environment
dev-setup: dev-db-up
	@echo "Development environment ready!"
	@echo "Run 'make run' to start the application"

dev-clean: dev-db-down
	@echo "Development environment cleaned up"

# Database Migration Commands
migrate-up:
	@echo "Running database migrations..."
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command up; \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command up; \
	fi

migrate-down:
	@echo "Rolling back last migration..."
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command down; \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command down; \
	fi

migrate-reset:
	@echo "Resetting all migrations..."
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command reset; \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command reset; \
	fi

migrate-status:
	@echo "Checking migration status..."
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command status; \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command status; \
	fi

migrate-version:
	@echo "Getting migration version..."
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command version; \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command version; \
	fi

migrate-create:
	@echo "Creating new migration..."
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command create -name $(NAME); \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command create -name $(NAME); \
	fi

migrate-fix:
	@echo "Fixing migration versioning..."
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command fix; \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command fix; \
	fi

migrate-validate:
	@echo "Validating migration files..."
	@if [ -f $(BUILD_DIR)/$(MIGRATE_BINARY) ]; then \
		$(BUILD_DIR)/$(MIGRATE_BINARY) -command validate; \
	else \
		$(GOCMD) run $(MIGRATE_PATH) -command validate; \
	fi

# Elasticsearch Management Commands
es-status:
	@echo "Checking Elasticsearch status..."
	@if [ -f $(BUILD_DIR)/$(ES_MANAGER_BINARY) ]; then \
		$(BUILD_DIR)/$(ES_MANAGER_BINARY) -command status; \
	else \
		$(GOCMD) run $(ES_MANAGER_PATH) -command status; \
	fi

es-init:
	@echo "Initializing Elasticsearch indexes..."
	@if [ -f $(BUILD_DIR)/$(ES_MANAGER_BINARY) ]; then \
		$(BUILD_DIR)/$(ES_MANAGER_BINARY) -command init; \
	else \
		$(GOCMD) run $(ES_MANAGER_PATH) -command init; \
	fi

es-recreate:
	@echo "Recreating Elasticsearch indexes..."
	@if [ -f $(BUILD_DIR)/$(ES_MANAGER_BINARY) ]; then \
		$(BUILD_DIR)/$(ES_MANAGER_BINARY) -command recreate; \
	else \
		$(GOCMD) run $(ES_MANAGER_PATH) -command recreate; \
	fi

es-health:
	@echo "Checking Elasticsearch health..."
	@if [ -f $(BUILD_DIR)/$(ES_MANAGER_BINARY) ]; then \
		$(BUILD_DIR)/$(ES_MANAGER_BINARY) -command health; \
	else \
		$(GOCMD) run $(ES_MANAGER_PATH) -command health; \
	fi

# Data Sync Commands
sync-data:
	@echo "Syncing data to Elasticsearch..."
	@if [ -f $(BUILD_DIR)/$(DATA_SYNC_BINARY) ]; then \
		$(BUILD_DIR)/$(DATA_SYNC_BINARY); \
	else \
		$(GOCMD) run $(DATA_SYNC_PATH); \
	fi

sync-data-dry:
	@echo "Dry run data sync..."
	@if [ -f $(BUILD_DIR)/$(DATA_SYNC_BINARY) ]; then \
		$(BUILD_DIR)/$(DATA_SYNC_BINARY) -dry-run; \
	else \
		$(GOCMD) run $(DATA_SYNC_PATH) -dry-run; \
	fi


# Generate Swagger documentation
gen-swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o internal/docs; \
		echo "Swagger documentation generated in internal/docs/"; \
	else \
		echo "Swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Generate Wire dependency injection
gen-wire:
	@echo "Generating Wire dependency injection..."
	@if command -v wire > /dev/null; then \
		cd cmd/server && wire; \
		echo "Wire files generated"; \
	else \
		echo "Wire not installed. Install with: go install github.com/google/wire/cmd/wire@latest"; \
	fi

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/cosmtrek/air@latest
	$(GOGET) github.com/swaggo/swag/cmd/swag@latest
	$(GOGET) github.com/google/wire/cmd/wire@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed"

# Check Go version
check-go-version:
	@echo "Checking Go version..."
	@$(GOCMD) version
	@echo "Required Go version: $(GO_VERSION)"

# Setup development environment
setup: check-go-version install-tools deps
	@echo "Development environment setup completed"

# Show help
help:
	@echo "Available commands:"
	@echo ""
	@echo "Build and Run:"
	@echo "  build          - Build the application"
	@echo "  build-importer - Build the importer tool"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Run the application locally"
	@echo "  run-dev        - Run with hot reload (requires air)"
	@echo "  run-importer   - Run the importer tool"
	@echo ""
	@echo "Development Environment:"
	@echo "  dev-setup      - Setup development environment (start DB)"
	@echo "  dev-db-up      - Start PostgreSQL database only"
	@echo "  dev-db-down    - Stop PostgreSQL database"
	@echo "  dev-db-logs    - Show PostgreSQL logs"
	@echo "  dev-db-reset   - Reset PostgreSQL database (WARNING: deletes data)"
	@echo "  dev-clean      - Clean development environment"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-import    - Run import functionality tests"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo ""
	@echo "Database:"
	@echo "  migrate-up      - Run all pending migrations"
	@echo "  migrate-down    - Roll back the last migration"
	@echo "  migrate-reset   - Roll back all migrations"
	@echo "  migrate-status  - Show migration status"
	@echo "  migrate-version - Show current migration version"
	@echo "  migrate-create  - Create new migration (use NAME=migration_name)"
	@echo "  migrate-fix     - Fix migration versioning issues"
	@echo "  migrate-validate - Validate migration files"
	@echo ""
	@echo "Elasticsearch:"
	@echo "  es-status       - Check Elasticsearch status"
	@echo "  es-init         - Initialize Elasticsearch indexes"
	@echo "  es-recreate     - Recreate Elasticsearch indexes"
	@echo "  es-health       - Check Elasticsearch health"
	@echo ""
	@echo "Data Sync:"
	@echo "  sync-data       - Sync database data to Elasticsearch"
	@echo "  sync-data-dry   - Dry run data sync (no actual sync)"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up   - Start PostgreSQL with docker-compose"
	@echo "  docker-compose-down - Stop PostgreSQL"
	@echo "  docker-compose-up-migrate - Start PostgreSQL and run migrations"
	@echo "  docker-build-migrate - Build migration Docker image"
	@echo "  docker-migrate-up - Run migrations in Docker"
	@echo "  docker-migrate-down - Roll back migrations in Docker"
	@echo ""
	@echo "Documentation:"
	@echo "  gen-swagger    - Generate Swagger documentation"
	@echo "  gen-wire       - Generate Wire dependency injection"
	@echo ""
	@echo "Setup:"
	@echo "  deps           - Download dependencies"
	@echo "  fix-deps       - Fix dependency issues (clean cache and reinstall)"
	@echo "  install-tools  - Install development tools"
	@echo "  setup          - Setup development environment"
	@echo "  check-go-version - Check Go version"
	@echo ""
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Quick Start for Local Development:"
	@echo "  make dev-setup    # Start database"
	@echo "  make run-dev      # Run app with hot reload"
