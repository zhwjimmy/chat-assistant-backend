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
- **Docker Support**: Multi-stage Docker build with PostgreSQL service
- **Database Migrations**: Goose-based database migration system
- **Dependency Injection**: Google Wire for compile-time DI
- **API Documentation**: OpenAPI 3.0 specification
- **Code Quality**: golangci-lint integration
- **Testing**: Unit test support

## Prerequisites

- Go 1.23.1 (or latest stable version)
- PostgreSQL 15+
- Docker & Docker Compose (for PostgreSQL service)

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
│   ├── server/           # Application entry point
│   ├── importer/         # Data import tool
│   └── migrate/          # Database migration tool
├── internal/
│   ├── app/              # Application bootstrap
│   ├── config/           # Configuration management
│   ├── server/           # HTTP server setup
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business logic
│   ├── repositories/     # Data access layer
│   ├── models/           # Data models
│   ├── migrations/       # Database migration files
│   ├── importer/         # Data import functionality
│   ├── i18n/             # Internationalization
│   ├── docs/             # API documentation
│   ├── logger/           # Logging setup
│   └── errors/           # Error handling
├── api/                  # OpenAPI specifications
├── config/               # Configuration files
├── scripts/              # Database scripts and sample data
├── test/                 # Test files
├── Dockerfile
├── Dockerfile.migrate    # Migration service Dockerfile
├── docker-compose.yaml
├── goose.yaml           # Goose migration configuration
├── Makefile
└── README.md
```

### Available Commands

```bash
# Build and run
make build              # Build the application
make run                # Run locally
make run-dev            # Run with hot reload

# Database migrations
make migrate-up         # Run all pending migrations
make migrate-down       # Roll back the last migration
make migrate-reset      # Roll back all migrations
make migrate-status     # Show migration status
make migrate-version    # Show current migration version
make migrate-create     # Create new migration (use NAME=migration_name)
make migrate-fix        # Fix migration versioning issues
make migrate-validate   # Validate migration files

# Testing
make test               # Run tests
make test-coverage      # Run tests with coverage

# Code quality
make lint               # Run linter
make fmt                # Format code
make vet                # Vet code

# Documentation
make gen-swagger        # Generate Swagger docs
make gen-wire           # Generate Wire DI

# Docker
make docker-build       # Build Docker image
make docker-run         # Run Docker container
make docker-compose-up  # Start PostgreSQL with docker-compose
make docker-compose-down # Stop PostgreSQL
make docker-compose-up-migrate # Start PostgreSQL and run migrations
make docker-build-migrate # Build migration Docker image
make docker-migrate-up  # Run migrations in Docker
make docker-migrate-down # Roll back migrations in Docker

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


# API documentation
go install github.com/swaggo/swag/cmd/swag@latest

# Dependency injection
go install github.com/google/wire/cmd/wire@latest

# Linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Database migrations (optional - migrations can be run via make commands)
go install github.com/pressly/goose/v3/cmd/goose@latest
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

### Docker Compose (PostgreSQL Service)

```bash
# Start PostgreSQL service
make docker-compose-up

# Stop PostgreSQL service
make docker-compose-down
```

The docker-compose setup includes:
- PostgreSQL database for local development
- Optional migration service for automatic database setup

## Database Migrations

The application uses [Goose](https://github.com/pressly/goose) for database migrations. Migrations are automatically run when the application starts, but can also be managed manually.

### Migration Files

Migration files are located in `internal/migrations/` and follow the naming convention:
- `{version}_{description}.up.sql` - Forward migration
- `{version}_{description}.down.sql` - Rollback migration

### Automatic Migrations

Migrations run automatically when the application starts. This ensures the database schema is always up-to-date.

### Manual Migration Management

```bash
# Run all pending migrations
make migrate-up

# Roll back the last migration
make migrate-down

# Reset all migrations (WARNING: This will drop all data)
make migrate-reset

# Check migration status
make migrate-status

# Get current migration version
make migrate-version

# Create a new migration
make migrate-create NAME=add_new_table

# Fix migration versioning issues
make migrate-fix

# Validate migration files
make migrate-validate
```

### Docker-based Migrations

```bash
# Start PostgreSQL and run migrations automatically
make docker-compose-up-migrate

# Run migrations in Docker container
make docker-migrate-up

# Roll back migrations in Docker container
make docker-migrate-down
```

### Migration Configuration

Migration settings can be configured in `goose.yaml`:

```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: chat_assistant
  sslmode: disable

migrations:
  dir: internal/migrations
  table: goose_db_version
  allow_missing: false
  allow_out_of_order: false
```

### Creating New Migrations

1. Create a new migration:
   ```bash
   make migrate-create NAME=add_user_preferences
   ```

2. Edit the generated files in `internal/migrations/`:
   - `{timestamp}_add_user_preferences.up.sql` - Add your schema changes
   - `{timestamp}_add_user_preferences.down.sql` - Add rollback changes

3. Test the migration:
   ```bash
   make migrate-up
   make migrate-down
   make migrate-up
   ```

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

3. **Database Errors**
   - Check database connection
   - Ensure database user has proper permissions
   - Verify database schema is properly initialized
   - Run `make migrate-status` to check migration status
   - Run `make migrate-up` to apply pending migrations

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

### v1.1.0
- Added Goose-based database migration system
- Automatic migration execution on application startup
- Migration management tools and commands
- Docker support for migrations
- Enhanced documentation for database management

### v1.0.0
- Initial release
- Basic API structure
- Database integration
- Docker support
- Documentation
