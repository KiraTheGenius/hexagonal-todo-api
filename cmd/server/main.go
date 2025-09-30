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

	"taskflow/adapter/cache"
	router "taskflow/adapter/http"
	"taskflow/adapter/http/handlers"
	"taskflow/adapter/repository"
	"taskflow/adapter/storage"
	"taskflow/adapter/streaming"
	"taskflow/internal/domain/file"
	"taskflow/internal/domain/todo"
	"taskflow/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewGormConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := repository.RunGormMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	todoRepo := repository.NewTodoRepository(db)

	s3Client := storage.NewS3Client(cfg.S3Config)
	fileStorage := storage.NewS3Storage(s3Client, cfg.S3Config.Bucket)
	fileRepo := repository.NewFileRepository(db) // Assuming we have a file repository

	redisClient := streaming.NewRedisClient(cfg.RedisURL)
	messaging := streaming.NewRedisMessaging(redisClient)
	cache := cache.NewRedisCache(redisClient)

	todoService := todo.NewTodoService(todoRepo, messaging, cache)
	fileService := file.NewFileService(fileRepo, fileStorage)

	todoHandler := handlers.NewTodoHandler(todoService)
	fileHandler := handlers.NewFileHandler(fileService)

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
