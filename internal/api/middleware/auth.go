package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// authMiddleware verifies that a valid user ID is present in the UserId header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("User-Id")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: Missing or invalid user ID",
			})
			return
		}

		userId := token

		c.Set("userId", userId)
		c.Next()
	}
}
