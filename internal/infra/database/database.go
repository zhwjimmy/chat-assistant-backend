package database

import (
	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/migrations"

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
