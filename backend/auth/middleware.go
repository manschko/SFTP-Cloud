package auth

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

// Middleware creates a Gin middleware for JWT authentication
func Middleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := authHeader[7:]

		// Parse and validate token
		claims, err := ValidateToken(tokenString, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Store user info in context
		c.Set("username", claims.Username)
		c.Next()
	}
}