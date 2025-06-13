package middleware

import (
	"net/http"
	"strings"

	"gocheck/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bearer TOKEN'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// ✅ Store both userID and userRole in context
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role) // this is what was missing

		c.Next()
	}
}

func RoleAuthorization(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleAny, exists := c.Get("userRole") // must match key used in AuthMiddleware
		role, ok := roleAny.(string)
		if !exists || !ok || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
