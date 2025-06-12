package routes

import (
	"gocheck/controllers"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(r *gin.Engine, db *gorm.DB) {
	bookController := controllers.NewBookController(db)

	bookRoutes := r.Group("/books")
	{
		bookRoutes.POST("/", bookController.CreateBook)      // Create a new book
		bookRoutes.GET("/", bookController.GetAllBooks)      // Get all books
		bookRoutes.GET("/:id", bookController.GetBookByID)   // Get a book by ID
		bookRoutes.PUT("/:id", bookController.UpdateBook)    // Update a book by ID
		bookRoutes.DELETE("/:id", bookController.DeleteBook) // Delete a book by ID
	}
}
