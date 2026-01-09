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

	// Retry connection with backoff (postgres may not be ready yet)
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if err = DB.PingContext(ctx); err == nil {
			log.Println("Database connected successfully")
			return nil
		}
		log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
