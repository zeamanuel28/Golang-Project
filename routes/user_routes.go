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
		// Public routes for demonstration (e.g., user creation doesn't require auth)
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.GET("/", userController.GetAllUsers)
		userRoutes.GET("/login", userController.Login)

		// Authenticated routes
		authenticatedRoutes := userRoutes.Group("/")
		authenticatedRoutes.Use(middleware.AuthMiddleware()) // Apply authentication middleware
		{

			authenticatedRoutes.PUT("/:id", userController.UpdateUser)
			authenticatedRoutes.DELETE("/:id", userController.DeleteUser)
		}
	}
}
