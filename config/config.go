package config

import (
	"fmt"
	"log"
	"os"
	"strconv" // Needed for parsing integers

	"github.com/joho/godotenv"
)

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	User     string
	Password string
	Name     string
	Port     int
	DSN      string // Data Source Name for GORM
}

// AppConfiguration holds all application-wide configuration
type AppConfiguration struct {
	Port      string // Changed to string to directly use os.Getenv result for router.Run
	JWTSecret string
	Database  DatabaseConfig
}

// AppConfig is the global instance of your application's configuration
var AppConfig AppConfiguration

// LoadConfig reads configuration from environment variables and populates AppConfig
func LoadConfig() error {
	// Load .env file. This will not error if the file doesn't exist,
	// allowing system environment variables to be used directly.
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// --- Load Application Port ---
	AppConfig.Port = os.Getenv("PORT")
	if AppConfig.Port == "" {
		// Provide a default port or make it a fatal error if port is essential
		// For now, let's default to "8080" if not set, but log a warning.
		log.Println("Warning: PORT environment variable not set, defaulting to 8080.")
		AppConfig.Port = "8080"
		// If you want to enforce PORT, use:
		// return fmt.Errorf("PORT environment variable not set")
	}

	// --- Load JWT Secret ---
	AppConfig.JWTSecret = os.Getenv("JWT_SECRET")
	if AppConfig.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable not set. Please set a strong, random secret")
	}

	// --- Load Database Configuration ---
	AppConfig.Database.Driver = os.Getenv("DB_DRIVER")
	if AppConfig.Database.Driver == "" {
		return fmt.Errorf("DB_DRIVER environment variable not set. Please specify 'postgres', 'sqlite3', etc.")
	}

	AppConfig.Database.Host = os.Getenv("DB_HOST")
	if AppConfig.Database.Host == "" {
		return fmt.Errorf("DB_HOST environment variable not set")
	}

	AppConfig.Database.User = os.Getenv("DB_USER")
	if AppConfig.Database.User == "" {
		return fmt.Errorf("DB_USER environment variable not set")
	}

	AppConfig.Database.Password = os.Getenv("DB_PASSWORD")
	// Password can be empty, so no error check here.

	AppConfig.Database.Name = os.Getenv("DB_NAME")
	if AppConfig.Database.Name == "" {
		return fmt.Errorf("DB_NAME environment variable not set")
	}

	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		return fmt.Errorf("DB_PORT environment variable not set")
	}
	port, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return fmt.Errorf("error parsing DB_PORT '%s': %w", dbPortStr, err)
	}
	AppConfig.Database.Port = port

	// Construct the DSN (Data Source Name) based on the selected driver
	switch AppConfig.Database.Driver {
	case "postgres":
		AppConfig.Database.DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			AppConfig.Database.Host,
			AppConfig.Database.User,
			AppConfig.Database.Password,
			AppConfig.Database.Name,
			AppConfig.Database.Port,
		)
	// You can add other database drivers here if needed, e.g., "sqlite3":
	// case "sqlite3":
	// 	AppConfig.Database.DSN = AppConfig.Database.Name // For SQLite, DSN is often just the file path
	default:
		return fmt.Errorf("unsupported database driver configured: %s", AppConfig.Database.Driver)
	}

	log.Println("Configuration loaded successfully.")
	return nil
}
