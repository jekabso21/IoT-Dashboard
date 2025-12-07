package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jekabso21/IoT-Dashboard/backend/internal/config"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.User,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.DBName,
		config.AppConfig.Database.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := DB.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func SetUserContext(ctx context.Context, userID string) (context.Context, error) {
	_, err := DB.ExecContext(ctx, "SET LOCAL app.current_user_id = $1", userID)
	if err != nil {
		return ctx, fmt.Errorf("failed to set user context: %w", err)
	}
	return ctx, nil
}

func WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
