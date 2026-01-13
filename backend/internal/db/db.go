package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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

			// Run migrations automatically
			if err := runMigrations(ctx); err != nil {
				return fmt.Errorf("failed to run migrations: %w", err)
			}

			return nil
		}
		log.Printf("Database connection attempt %d/%d failed: %v", i+1, maxRetries, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// runMigrations executes all SQL migration files in order
func runMigrations(ctx context.Context) error {
	// Get all migration files
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Filter and sort SQL files
	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Execute each migration
	for _, filename := range sqlFiles {
		content, err := fs.ReadFile(migrationsFS, "migrations/"+filename)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		log.Printf("Running migration: %s", filename)
		if _, err := DB.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}
	}

	log.Printf("Migrations completed successfully (%d files)", len(sqlFiles))
	return nil
}

// Close closes the database connection pool
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
