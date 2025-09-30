package http

import (
	"taskflow/adapter/http/handlers"
	"taskflow/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(todoHandler *handlers.TodoHandler, fileHandler *handlers.FileHandler) *gin.Engine {
	r := gin.New()

	// Add middleware
	r.Use(middleware.RequestLogger(middleware.DefaultRequestLoggerConfig()))
	r.Use(middleware.CORS(middleware.DefaultCORSConfig()))
	r.Use(middleware.Recovery(middleware.DefaultRecoveryConfig()))
	r.Use(middleware.Timeout(middleware.DefaultTimeoutConfig()))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// File upload
	r.POST("/upload", fileHandler.UploadFile)

	// Todo endpoints
	todoGroup := r.Group("/todo")
	{
		todoGroup.POST("", todoHandler.CreateTodo)
		todoGroup.GET("/:id", todoHandler.GetTodo)
		todoGroup.GET("", todoHandler.ListTodos)
		todoGroup.PUT("/:id", todoHandler.UpdateTodo)
		todoGroup.DELETE("/:id", todoHandler.DeleteTodo)
	}

	return r
}
