package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
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

// dbLogger holds the logger for database operations
var dbLogger *slog.Logger

// SetLogger sets the logger for database operations
func SetLogger(log *slog.Logger) {
	dbLogger = log
}

// logInfo logs at INFO level, falling back to standard log if no logger set
func logInfo(msg string, args ...any) {
	if dbLogger != nil {
		dbLogger.Info(msg, args...)
	} else {
		log.Println(msg)
	}
}

// logError logs at ERROR level, falling back to standard log if no logger set
func logError(msg string, args ...any) {
	if dbLogger != nil {
		dbLogger.Error(msg, args...)
	} else {
		log.Println("ERROR:", msg)
	}
}

// logWarn logs at WARN level, falling back to standard log if no logger set
func logWarn(msg string, args ...any) {
	if dbLogger != nil {
		dbLogger.Warn(msg, args...)
	} else {
		log.Println("WARN:", msg)
	}
}

func Connect(ctx context.Context) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logWarn("database_skip", slog.String("reason", "DATABASE_URL not set"))
		return nil
	}

	logInfo("database_connecting")

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		logError("database_open_failed", slog.String("error", err.Error()))
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
		logInfo("database_connection_attempt",
			slog.Int("attempt", i+1),
			slog.Int("max_attempts", maxRetries),
		)

		if err = DB.PingContext(ctx); err == nil {
			logInfo("database_connected",
				slog.Int("max_open_conns", defaultMaxOpenConns),
				slog.Int("max_idle_conns", defaultMaxIdleConns),
				slog.Duration("conn_max_lifetime", defaultConnMaxLifetime),
				slog.Duration("conn_max_idle_time", defaultConnMaxIdleTime),
			)

			// Run migrations automatically
			if err := runMigrations(ctx); err != nil {
				logError("database_migrations_failed", slog.String("error", err.Error()))
				return fmt.Errorf("failed to run migrations: %w", err)
			}

			return nil
		}

		logWarn("database_connection_retry",
			slog.Int("attempt", i+1),
			slog.Int("max_attempts", maxRetries),
			slog.String("error", err.Error()),
			slog.Duration("retry_delay", time.Duration(i+1)*time.Second),
		)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	logError("database_connection_failed",
		slog.Int("attempts", maxRetries),
		slog.String("error", err.Error()),
	)
	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// runMigrations executes all SQL migration files in order
func runMigrations(ctx context.Context) error {
	logInfo("migrations_starting")

	// Get all migration files
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		logError("migrations_read_dir_failed", slog.String("error", err.Error()))
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

	logInfo("migrations_found", slog.Int("count", len(sqlFiles)))

	// Execute each migration
	for _, filename := range sqlFiles {
		content, err := fs.ReadFile(migrationsFS, "migrations/"+filename)
		if err != nil {
			logError("migration_read_failed",
				slog.String("file", filename),
				slog.String("error", err.Error()),
			)
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		logInfo("migration_executing", slog.String("file", filename))
		start := time.Now()

		if _, err := DB.ExecContext(ctx, string(content)); err != nil {
			logError("migration_failed",
				slog.String("file", filename),
				slog.String("error", err.Error()),
				slog.Duration("duration", time.Since(start)),
			)
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		logInfo("migration_completed",
			slog.String("file", filename),
			slog.Duration("duration", time.Since(start)),
		)
	}

	logInfo("migrations_completed", slog.Int("count", len(sqlFiles)))
	return nil
}

// Close closes the database connection pool
func Close() error {
	if DB != nil {
		logInfo("database_closing")
		err := DB.Close()
		if err != nil {
			logError("database_close_failed", slog.String("error", err.Error()))
		} else {
			logInfo("database_closed")
		}
		return err
	}
	return nil
}
