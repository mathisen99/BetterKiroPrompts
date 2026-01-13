package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connection pool configuration
const (
	defaultMaxOpenConns    = 25
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = 5 * time.Minute
	defaultConnMaxIdleTime = 1 * time.Minute
)

var DB *sql.DB

func Connect(ctx context.Context) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL not set, skipping database connection")
		return nil
	}

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return err
	}

	// Configure connection pool for concurrent access
	DB.SetMaxOpenConns(defaultMaxOpenConns)
	DB.SetMaxIdleConns(defaultMaxIdleConns)
	DB.SetConnMaxLifetime(defaultConnMaxLifetime)
	DB.SetConnMaxIdleTime(defaultConnMaxIdleTime)

	// Retry connection with backoff (postgres may not be ready yet)
	maxRetries := 5
	for i := range maxRetries {
		if err = DB.PingContext(ctx); err == nil {
			log.Printf("Database connected successfully (pool: max=%d, idle=%d)",
				defaultMaxOpenConns, defaultMaxIdleConns)
			return nil
		}
		log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// Close closes the database connection pool
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
