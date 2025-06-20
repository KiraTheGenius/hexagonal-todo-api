package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	S3Config    S3Config
	Environment string
	LogLevel    string
}

type S3Config struct {
	Region          string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "root:password@tcp(localhost:3306)/todoservice?parseTime=true"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		S3Config: S3Config{
			Region:          getEnv("AWS_REGION", "us-east-1"),
			Bucket:          getEnv("S3_BUCKET", "todo-files"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", "test"),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", "test"),
			Endpoint:        getEnv("S3_ENDPOINT", "http://localhost:4566"),
		},
	}

	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("Configuration validation failed: %v", err))
	}

	return cfg
}

func (c *Config) Validate() error {
	// Validate port
	if port, err := strconv.Atoi(c.Port); err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %s", c.Port)
	}

	// Validate database URL
	if c.DatabaseURL == "" {
		return fmt.Errorf("database URL is required")
	}

	// Validate Redis URL
	if c.RedisURL == "" {
		return fmt.Errorf("Redis URL is required")
	}

	// Validate S3 configuration
	if c.S3Config.Bucket == "" {
		return fmt.Errorf("S3 bucket is required")
	}

	// Validate environment
	validEnvironments := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
	}
	if !validEnvironments[c.Environment] {
		return fmt.Errorf("invalid environment: %s", c.Environment)
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s", c.LogLevel)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
