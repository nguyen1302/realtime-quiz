package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"realtime-quiz/internal/bootstrap"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Real-Time Quiz Server")

	// Determine config path from environment or use default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	// Load configuration from YAML
	cfg, err := bootstrap.LoadConfig(configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err, "path", configPath)
		os.Exit(1)
	}

	slog.Info("Configuration loaded", "path", configPath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize PostgreSQL
	pgPool, err := bootstrap.NewPostgresPool(ctx, &cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer bootstrap.ClosePostgres(pgPool)

	// Initialize Redis
	redisClient, err := bootstrap.NewRedisClient(ctx, &cfg.Redis)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer bootstrap.CloseRedis(redisClient)

	slog.Info("All connections established successfully! ✅")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")
}
