package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeoutConfig represents timeout middleware configuration
type TimeoutConfig struct {
	Timeout time.Duration
	Logger  *slog.Logger
}

// DefaultTimeoutConfig returns a default timeout configuration
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		Timeout: 30 * time.Second,
		Logger:  slog.Default(),
	}
}

// Timeout returns a timeout middleware
func Timeout(config TimeoutConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.Timeout)
		defer cancel()

		// Replace request context
		c.Request = c.Request.WithContext(ctx)

		// Channel to signal completion
		done := make(chan struct{})

		// Start goroutine to handle the request
		go func() {
			defer func() {
				if r := recover(); r != nil {
					config.Logger.Error("panic in timeout middleware", "error", r)
				}
				close(done)
			}()
			c.Next()
		}()

		// Wait for completion or timeout
		select {
		case <-done:
			// Request completed normally
			return
		case <-ctx.Done():
			// Timeout occurred
			config.Logger.Warn("request timeout",
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"timeout", config.Timeout,
			)

			// Check if response was already sent
			if !c.Writer.Written() {
				c.JSON(http.StatusRequestTimeout, gin.H{
					"error":   "Request timeout",
					"timeout": config.Timeout.String(),
				})
			}

			c.Abort()
		}
	}
}
