package main

import (
	"gocheck/config"
	"gocheck/database"
	"gocheck/routes"

	_ "gocheck/docs"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Create a single Gin router instance
	router := gin.Default()

	// Register your application routes on this router
	routes.SetupUserRoutes(router, db)
	routes.RegisterBookRoutes(router, db)

	// Register swagger handler on the same router
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server on the configured port
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
