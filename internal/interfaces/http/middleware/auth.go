package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If no token is configured, skip authentication
		if token == "" {
			c.Next()
			return
		}

		// Get token from header
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "Authorization token is required",
				"data":    nil,
			})
			return
		}

		// Validate token
		if authToken != token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    http.StatusUnauthorized,
				"message": "Invalid authorization token",
				"data":    nil,
			})
			return
		}

		c.Next()
	}
}
