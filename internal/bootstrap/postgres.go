package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nguyen1302/realtime-quiz/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormDB creates a new GORM database connection
func NewGormDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Info("PostgreSQL connected successfully (GORM)",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName,
	)

	return db, nil
}

// CloseGormDB closes the GORM database connection
func CloseGormDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			slog.Error("failed to get sql.DB for closing", "error", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
		} else {
			slog.Info("Database connection closed")
		}
	}
}
