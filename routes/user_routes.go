package routes

import (
	"gocheck/controllers"
	"gocheck/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupUserRoutes sets up the routes for user-related operations
func SetupUserRoutes(router *gin.Engine, db *gorm.DB) {
	userController := controllers.NewUserController(db)

	userRoutes := router.Group("/users")
	{
		// Public routes (no auth required)
		userRoutes.POST("/", userController.CreateUser)    // Register user
		userRoutes.POST("/login", userController.Login)    // Login (POST)
		userRoutes.GET("/:id", userController.GetUserByID) // Get user by ID
		userRoutes.GET("/", userController.GetAllUsers)    // Get all users

		// Routes that require authentication
		authenticatedRoutes := userRoutes.Group("/")
		authenticatedRoutes.Use(middleware.AuthMiddleware())
		{
			authenticatedRoutes.PUT("/:id", userController.UpdateUser) // Update user info

			// Admin-only routes nested under authenticated routes
			adminRoutes := authenticatedRoutes.Group("/admin")
			adminRoutes.Use(middleware.RoleAuthorization("admin"))
			{
				adminRoutes.DELETE("/:id", userController.DeleteUser) // Delete user (admin only)
			}
		}
	}
}
