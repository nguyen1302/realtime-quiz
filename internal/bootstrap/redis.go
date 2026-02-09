package bootstrap

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nguyen1302/realtime-quiz/internal/config"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates a new Redis client
func NewRedisClient(ctx context.Context, cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: 10,
	})

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	slog.Info("Redis connected successfully",
		"host", cfg.Host,
		"port", cfg.Port,
	)

	return client, nil
}

// CloseRedis closes the Redis client connection
func CloseRedis(client *redis.Client) {
	if client != nil {
		if err := client.Close(); err != nil {
			slog.Error("failed to close redis connection", "error", err)
		} else {
			slog.Info("Redis connection closed")
		}
	}
}
