package routes

import (
	"gocheck/controllers"
	"gocheck/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	// Public endpoints (Swagger will pick these up)
	router.POST("/users", userController.CreateUser)     // Register user
	router.POST("/login", userController.Login)          // Login
	router.GET("/users/:id", userController.GetUserByID) // Get user by ID
	router.GET("/users", userController.GetAllUsers)     // Get all users

	// Protected routes (also visible to Swagger)
	router.PUT("/users/:id", middleware.AuthMiddleware(), userController.UpdateUser)

	// Admin-only route
	router.DELETE("/users/admin/:id",
		middleware.AuthMiddleware(),
		middleware.RoleAuthorization("admin"),
		userController.DeleteUser,
	)
}
