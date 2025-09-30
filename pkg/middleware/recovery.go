package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
)

// RecoveryConfig represents recovery middleware configuration
type RecoveryConfig struct {
	Logger          *slog.Logger
	StackAll        bool
	StackSize       int
	DisableStackAll bool
}

// DefaultRecoveryConfig returns a default recovery configuration
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		Logger:          slog.Default(),
		StackAll:        false,
		StackSize:       1024 * 8,
		DisableStackAll: false,
	}
}

// Recovery returns a panic recovery middleware
func Recovery(config RecoveryConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// Get request info
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// Build log fields
				fields := []interface{}{
					"error", err,
					"request", string(httpRequest),
				}

				// Add stack trace if enabled
				if !config.DisableStackAll {
					stack := debug.Stack()
					if config.StackAll {
						fields = append(fields, "stack", string(stack))
					} else {
						fields = append(fields, "stack", string(stack[:config.StackSize]))
					}
				}

				// Log the panic
				if brokenPipe {
					config.Logger.Error("broken pipe", fields...)
				} else {
					config.Logger.Error("panic recovered", fields...)
				}

				// If the connection is dead, we can't write a status to it
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				// Return 500 error
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
