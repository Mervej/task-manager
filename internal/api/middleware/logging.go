package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		log.Printf("Started %s %s", method, path)

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("Completed %s %s with status %d in %v", method, path, statusCode, latency)
	}
}
