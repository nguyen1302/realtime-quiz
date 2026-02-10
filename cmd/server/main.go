package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nguyen1302/realtime-quiz/internal/bootstrap"
	"github.com/nguyen1302/realtime-quiz/internal/config"
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
	cfg, err := config.Load(configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err, "path", configPath)
		os.Exit(1)
	}

	slog.Info("Configuration loaded", "path", configPath)

	// Initialize PostgreSQL (GORM)
	db, err := bootstrap.NewGormDB(&cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to PostgreSQL", "error", err)
		os.Exit(1)
	}
	defer bootstrap.CloseGormDB(db)

	// Initialize Redis
	ctx := context.Background()
	redisClient, err := bootstrap.NewRedisClient(ctx, &cfg.Redis)
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer bootstrap.CloseRedis(redisClient)

	// Initialize Router
	router := bootstrap.NewRouter(db, redisClient, cfg)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: router.Engine(),
	}

	// Start server in goroutine
	go func() {
		slog.Info("HTTP server starting", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	slog.Info("Server is running! ✅", "address", addr)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited properly")
}
