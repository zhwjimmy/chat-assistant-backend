# Chat Assistant Backend

A Golang backend service for the Chat Assistant application, built with Gin, GORM, and PostgreSQL.

## Features

- **RESTful API**: Clean and well-documented API endpoints
- **Database**: PostgreSQL with GORM ORM
- **Configuration**: YAML-based configuration with environment variable support
- **Logging**: Structured JSON logging with Zap
- **CORS**: Configurable CORS support
- **Internationalization**: Multi-language support (English/Chinese)
- **Graceful Shutdown**: Proper signal handling and graceful shutdown
- **Request ID**: Request tracing with unique IDs
- **Docker Support**: Multi-stage Docker build with docker-compose
- **Database Migrations**: Goose migration support
- **Dependency Injection**: Google Wire for compile-time DI
- **API Documentation**: OpenAPI 3.0 specification
- **Code Quality**: golangci-lint integration
- **Testing**: Unit test support

## Prerequisites

- Go 1.23.1 (or latest stable version)
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Go Version Management

If Go 1.23.1 is not available, you can use version managers:

```bash
# Using gvm
gvm install go1.23.1
gvm use go1.23.1

# Using asdf
asdf install golang 1.23.1
asdf global golang 1.23.1

# Using goenv
goenv install 1.23.1
goenv global 1.23.1
```

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd chat-assistant-backend
make setup
```

### 2. Environment Configuration

Copy the example environment file and configure:

```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Database Setup

Start PostgreSQL with Docker:

```bash
make docker-compose-up
```

Or run migrations manually:

```bash
make migrate-up
```

### 4. Run the Application

```bash
# Development mode with hot reload
make run-dev

# Or standard run
make run
```

The API will be available at `http://localhost:8080`

## Development

### Project Structure

```
chat-assistant-backend/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── app/              # Application bootstrap
│   ├── config/           # Configuration management
│   ├── server/           # HTTP server setup
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business logic
│   ├── repositories/     # Data access layer
│   ├── models/           # Data models
│   ├── migrations/       # Database migrations
│   ├── i18n/             # Internationalization
│   ├── docs/             # API documentation
│   ├── logger/           # Logging setup
│   └── errors/           # Error handling
├── api/                  # OpenAPI specifications
├── config/               # Configuration files
├── test/                 # Test files
├── Dockerfile
├── docker-compose.yaml
├── Makefile
└── README.md
```

### Available Commands

```bash
# Build and run
make build              # Build the application
make run                # Run locally
make run-dev            # Run with hot reload

# Testing
make test               # Run tests
make test-coverage      # Run tests with coverage

# Code quality
make lint               # Run linter
make fmt                # Format code
make vet                # Vet code

# Database
make migrate-up         # Run migrations up
make migrate-down       # Run migrations down

# Documentation
make gen-swagger        # Generate Swagger docs
make gen-wire           # Generate Wire DI

# Docker
make docker-build       # Build Docker image
make docker-run         # Run Docker container
make docker-compose-up  # Start with docker-compose
make docker-compose-down # Stop docker-compose

# Development tools
make install-tools      # Install development tools
make setup              # Setup development environment
```

### Required Tools Installation

Some commands require external tools. Install them with:

```bash
make install-tools
```

Or install manually:

```bash
# Hot reload
go install github.com/cosmtrek/air@latest

# Database migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# API documentation
go install github.com/swaggo/swag/cmd/swag@latest

# Dependency injection
go install github.com/google/wire/cmd/wire@latest

# Linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `chat_assistant` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `SERVER_HOST` | Server host | `0.0.0.0` |
| `SERVER_PORT` | Server port | `8080` |
| `ALLOWED_ORIGINS` | CORS allowed origins | `http://localhost:3000` |
| `LOG_LEVEL` | Log level | `info` |
| `LOG_FORMAT` | Log format | `json` |
| `DEFAULT_LANGUAGE` | Default language | `en` |
| `SHUTDOWN_TIMEOUT` | Shutdown timeout | `30s` |

### Configuration File

The application uses `config/config.yaml` for default configuration. Environment variables override file settings.

## API Documentation

### Health Check

```bash
GET /health
```

Returns service health status.

### API Endpoints

The API follows RESTful conventions with standard JSON responses:

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

Error responses:

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": "Additional details"
  }
}
```

## Docker Deployment

### Build and Run

```bash
# Build image
make docker-build

# Run container
make docker-run
```

### Docker Compose

```bash
# Start all services
make docker-compose-up

# Stop all services
make docker-compose-down
```

The docker-compose setup includes:
- PostgreSQL database
- Chat Assistant Backend
- Optional migration service

## Database Migrations

### Using Goose

```bash
# Run migrations up
make migrate-up

# Run migrations down
make migrate-down
```

### Migration Files

Migration files are located in `internal/migrations/` and follow the naming convention:
- `YYYYMMDD_HHMMSS_description.up.sql`
- `YYYYMMDD_HHMMSS_description.down.sql`

## Internationalization

The application supports multiple languages through the i18n system:

- English (`en`) - Default
- Chinese (`zh`)

Language files are located in `internal/i18n/locales/`.

## Logging

The application uses structured JSON logging with the following fields:
- `timestamp`: Log timestamp
- `level`: Log level (debug, info, warn, error)
- `msg`: Log message
- `request_id`: Request ID for tracing

## Error Handling

The application uses a unified error handling system with:
- Standardized error codes
- Localized error messages
- Detailed error information
- Proper HTTP status codes

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Ensure PostgreSQL is running
   - Check database credentials in `.env`
   - Verify database exists

2. **Port Already in Use**
   - Change `SERVER_PORT` in configuration
   - Kill existing processes on port 8080

3. **Migration Errors**
   - Check database connection
   - Verify migration files syntax
   - Ensure database user has proper permissions

4. **Build Failures**
   - Ensure Go version is 1.23.1 or compatible
   - Run `make deps` to update dependencies
   - Check for syntax errors

### Getting Help

- Check the logs for detailed error messages
- Review the configuration files
- Ensure all prerequisites are installed
- Check the API documentation

## Changelog

### v1.0.0
- Initial release
- Basic API structure
- Database integration
- Docker support
- Documentation
