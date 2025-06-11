package middleware

import (
	"net/http"
	"strings"

	"gocheck/utils" // Your utils package for JWT

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware to authenticate requests using JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort() // Stop processing this request
			return
		}

		// Check if the header starts with "Bearer "
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer TOKEN'"})
			c.Abort()
			return
		}

		tokenString := parts[1] // Extract the token string

		// Validate the token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// If the token is valid, set the user ID in the context
		// This makes the user ID available to subsequent handlers in the chain
		c.Set("userID", claims.UserID)

		// Continue to the next handler in the chain
		c.Next()
	}
}
