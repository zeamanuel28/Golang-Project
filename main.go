package main

import (
	"gocheck/config"
	"gocheck/database"
	"gocheck/routes"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Auto-migrate models
	database.Migrate(db)

	// Set Gin to release mode for production (or comment out for development)
	// gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.Default()

	// Pass the DB instance to routes
	routes.SetupUserRoutes(router, db)

	// Start the server
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
