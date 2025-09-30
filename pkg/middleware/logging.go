package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLoggerConfig represents request logger configuration
type RequestLoggerConfig struct {
	Logger     *slog.Logger
	SkipPaths  []string
	SkipFields []string
}

// DefaultRequestLoggerConfig returns a default request logger configuration
func DefaultRequestLoggerConfig() RequestLoggerConfig {
	return RequestLoggerConfig{
		Logger:     slog.Default(),
		SkipPaths:  []string{"/health", "/metrics"},
		SkipFields: []string{},
	}
}

// RequestLogger returns a request logging middleware
func RequestLogger(config RequestLoggerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip logging for certain paths
		if shouldSkipPath(c.Request.URL.Path, config.SkipPaths) {
			c.Next()
			return
		}

		// Generate request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response info
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		// Build log fields
		fields := []interface{}{
			"request_id", requestID,
			"method", method,
			"path", path,
			"status", status,
			"latency", latency,
			"client_ip", clientIP,
		}

		// Add query parameters if present
		if raw != "" {
			fields = append(fields, "query", raw)
		}

		// Add user agent
		if userAgent := c.GetHeader("User-Agent"); userAgent != "" {
			fields = append(fields, "user_agent", userAgent)
		}

		// Add referer
		if referer := c.GetHeader("Referer"); referer != "" {
			fields = append(fields, "referer", referer)
		}

		// Log the request
		config.Logger.Info("request completed", fields...)
	}
}

// shouldSkipPath checks if the path should be skipped
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}
