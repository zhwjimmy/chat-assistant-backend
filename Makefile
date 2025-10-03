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

.PHONY: all build clean test deps run docker-build docker-run migrate-up migrate-down gen-swagger gen-wire lint help dev-db-up dev-db-down dev-db-logs dev-db-reset dev-setup dev-clean

# Default target
all: deps build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PATH)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

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

# Docker compose up
docker-compose-up:
	@echo "Starting services with docker-compose..."
	docker-compose up -d

# Docker compose down
docker-compose-down:
	@echo "Stopping services with docker-compose..."
	docker-compose down

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

# Database migration up
migrate-up:
	@echo "Running database migrations up..."
	@if command -v goose > /dev/null; then \
		goose -dir internal/migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=chat_assistant sslmode=disable" up; \
	else \
		echo "Goose not installed. Install with: go install github.com/pressly/goose/v3/cmd/goose@latest"; \
	fi

# Database migration down
migrate-down:
	@echo "Running database migrations down..."
	@if command -v goose > /dev/null; then \
		goose -dir internal/migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=chat_assistant sslmode=disable" down; \
	else \
		echo "Goose not installed. Install with: go install github.com/pressly/goose/v3/cmd/goose@latest"; \
	fi

# Generate Swagger documentation
gen-swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/server/main.go -o internal/docs; \
		@echo "Swagger documentation generated in internal/docs/"; \
	else \
		echo "Swag not installed. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Generate Wire dependency injection
gen-wire:
	@echo "Generating Wire dependency injection..."
	@if command -v wire > /dev/null; then \
		cd internal && wire; \
		@echo "Wire files generated"; \
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
	$(GOGET) github.com/pressly/goose/v3/cmd/goose@latest
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
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Run the application locally"
	@echo "  run-dev        - Run with hot reload (requires air)"
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
	@echo ""
	@echo "Code Quality:"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo ""
	@echo "Database:"
	@echo "  migrate-up     - Run database migrations up"
	@echo "  migrate-down   - Run database migrations down"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  docker-compose-up   - Start all services with docker-compose"
	@echo "  docker-compose-down - Stop all services with docker-compose"
	@echo ""
	@echo "Documentation:"
	@echo "  gen-swagger    - Generate Swagger documentation"
	@echo "  gen-wire       - Generate Wire dependency injection"
	@echo ""
	@echo "Setup:"
	@echo "  deps           - Download dependencies"
	@echo "  install-tools  - Install development tools"
	@echo "  setup          - Setup development environment"
	@echo "  check-go-version - Check Go version"
	@echo ""
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Quick Start for Local Development:"
	@echo "  make dev-setup    # Start database"
	@echo "  make run-dev      # Run app with hot reload"
