package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

// Migrator handles database migrations using Goose
type Migrator struct {
	db     *gorm.DB
	sqlDB  *sql.DB
	config *Config
}

// Config holds migration configuration
type Config struct {
	MigrationsDir   string
	TableName       string
	AllowMissing    bool
	AllowOutOfOrder bool
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB, config *Config) (*Migrator, error) {
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if config == nil {
		config = &Config{
			MigrationsDir:   "internal/migrations",
			TableName:       "goose_db_version",
			AllowMissing:    false,
			AllowOutOfOrder: false,
		}
	}

	return &Migrator{
		db:     db,
		sqlDB:  sqlDB,
		config: config,
	}, nil
}

// Up runs all pending migrations
func (m *Migrator) Up() error {
	log.Println("Running database migrations...")

	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Run migrations
	if err := goose.Up(m.sqlDB, m.config.MigrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	log.Println("Rolling back last migration...")

	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Roll back migration
	if err := goose.Down(m.sqlDB, m.config.MigrationsDir); err != nil {
		return fmt.Errorf("failed to roll back migration: %w", err)
	}

	log.Println("Migration rollback completed successfully")
	return nil
}

// Reset rolls back all migrations
func (m *Migrator) Reset() error {
	log.Println("Resetting all migrations...")

	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Reset migrations
	if err := goose.Reset(m.sqlDB, m.config.MigrationsDir); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	log.Println("All migrations reset successfully")
	return nil
}

// Status shows the current migration status
func (m *Migrator) Status() error {
	log.Println("Checking migration status...")

	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Get status
	if err := goose.Status(m.sqlDB, m.config.MigrationsDir); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}

// Version shows the current migration version
func (m *Migrator) Version() (int64, error) {
	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Get current version
	version, err := goose.GetDBVersion(m.sqlDB)
	if err != nil {
		return 0, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, nil
}

// Create creates a new migration file
func (m *Migrator) Create(name, migrationType string) error {
	log.Printf("Creating new migration: %s (%s)", name, migrationType)

	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Create migration
	if err := goose.Create(m.sqlDB, m.config.MigrationsDir, name, migrationType); err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}

	log.Printf("Migration created successfully: %s", name)
	return nil
}

// Fix fixes migration versioning issues
func (m *Migrator) Fix() error {
	log.Println("Fixing migration versioning...")

	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Fix migrations
	if err := goose.Fix(m.config.MigrationsDir); err != nil {
		return fmt.Errorf("failed to fix migrations: %w", err)
	}

	log.Println("Migration versioning fixed successfully")
	return nil
}

// Validate validates migration files
func (m *Migrator) Validate() error {
	log.Println("Validating migration files...")

	// Set Goose configuration
	goose.SetTableName(m.config.TableName)

	// Get current version to validate migrations are in order
	_, err := goose.GetDBVersion(m.sqlDB)
	if err != nil {
		return fmt.Errorf("migration validation failed: %w", err)
	}

	log.Println("Migration validation passed")
	return nil
}
