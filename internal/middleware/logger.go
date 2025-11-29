package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs incoming requests with detailed information
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		// Process request
		ctx.Next()

		// Calculate request duration
		duration := time.Since(startTime)

		// Log request details
		log.Printf(
			"[%s] %s %s | Status: %d | Duration: %v | IP: %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			ctx.Request.Proto,
			ctx.Writer.Status(),
			duration,
			ctx.ClientIP(),
		)
	}
}
