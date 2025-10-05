package database

import (
	"github.com/google/wire"
)

// DatabaseSet provides all database-related dependencies
var DatabaseSet = wire.NewSet(
	InitializeDatabase,
	RunMigrations,
)
