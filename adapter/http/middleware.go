package http

import (
	"taskflow/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// CORS returns a CORS middleware using the shared middleware package
func CORS() gin.HandlerFunc {
	return middleware.CORS(middleware.DefaultCORSConfig())
}

// RequestLogger returns a request logger middleware using the shared middleware package
func RequestLogger() gin.HandlerFunc {
	return middleware.RequestLogger(middleware.DefaultRequestLoggerConfig())
}

// Recovery returns a recovery middleware using the shared middleware package
func Recovery() gin.HandlerFunc {
	return middleware.Recovery(middleware.DefaultRecoveryConfig())
}

// Timeout returns a timeout middleware using the shared middleware package
func Timeout() gin.HandlerFunc {
	return middleware.Timeout(middleware.DefaultTimeoutConfig())
}
