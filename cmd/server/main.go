package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskflow/internal/adapter/database"
	"taskflow/internal/adapter/storage"
	"taskflow/internal/adapter/streaming"
	router "taskflow/internal/interfaces/http"
	"taskflow/internal/interfaces/http/handlers"
	"taskflow/internal/service"
	"taskflow/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.NewGormConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.RunGormMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	todoRepo := database.NewTodoRepository(db)

	s3Client := storage.NewS3Client(cfg.S3Config)
	fileRepo := storage.NewFileRepository(s3Client, cfg.S3Config.Bucket)

	redisClient := streaming.NewRedisClient(cfg.RedisURL)
	streamRepo := streaming.NewStreamRepository(redisClient)

	todoSr := service.NewTodoService(todoRepo, streamRepo)
	fileSr := service.NewFileService(fileRepo)

	todoHandler := handlers.NewTodoHandler(todoSr)
	fileHandler := handlers.NewFileHandler(fileSr)

	gin.SetMode(gin.ReleaseMode)
	r := router.SetupRouter(todoHandler, fileHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Failed to start server:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
