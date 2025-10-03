package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"chat-assistant-backend/internal/config"
	"chat-assistant-backend/internal/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up, down, reset, status, version, create, fix, validate")
		name    = flag.String("name", "", "Migration name (for create command)")
		mtype   = flag.String("type", "sql", "Migration type: sql, go (for create command)")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create migrator
	migrator, err := migrations.NewMigrator(db, nil)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	// Execute command
	switch *command {
	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatalf("Failed to roll back migration: %v", err)
		}
	case "reset":
		if err := migrator.Reset(); err != nil {
			log.Fatalf("Failed to reset migrations: %v", err)
		}
	case "status":
		if err := migrator.Status(); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	case "version":
		version, err := migrator.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d\n", version)
	case "create":
		if *name == "" {
			log.Fatal("Migration name is required for create command")
		}
		if err := migrator.Create(*name, *mtype); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
	case "fix":
		if err := migrator.Fix(); err != nil {
			log.Fatalf("Failed to fix migrations: %v", err)
		}
	case "validate":
		if err := migrator.Validate(); err != nil {
			log.Fatalf("Failed to validate migrations: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", *command)
		fmt.Println("Available commands: up, down, reset, status, version, create, fix, validate")
		os.Exit(1)
	}
}
