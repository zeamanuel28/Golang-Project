package database

import (
	"fmt"
	"gocheck/config"
	"gocheck/models"
	"log"

	"gorm.io/driver/postgres" // Changed to PostgreSQL driver
	"gorm.io/gorm"
)

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dbDriver := config.AppConfig.Database.Driver
	dbDSN := config.AppConfig.Database.DSN

	switch dbDriver {
	case "postgres": // Handle PostgreSQL connection
		db, err = gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
		}
	// You can keep sqlite3 as an option if you want, but it's not needed for PostgreSQL
	// case "sqlite3":
	// 	db, err = gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	// 	if err != nil {
	// 		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	// 	}
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	log.Println("Database connection established successfully!")
	return db, nil
}

// Migrate runs database migrations (auto-migrate models)
func Migrate(db *gorm.DB) {
	log.Println("Running database migrations...")

	// Drop tables before migration ONLY in development/testing (optional)
	// ⚠️ WARNING: This will delete all data. Use carefully.
	// _ = db.Migrator().DropTable(&models.Book{}, &models.User{})

	// AutoMigrate should run in correct order: parent before child
	db.Migrator().DropTable(&models.Book{}, &models.User{})
	err := db.AutoMigrate(
		&models.User{},
		&models.Book{},
	)

	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database migrations completed.")
}
