// Package db provides database connectivity and logging wrappers.
package db

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"
	"time"

	"better-kiro-prompts/internal/logger"
)

// LoggingDB wraps sql.DB with query logging
type LoggingDB struct {
	db  *sql.DB
	log *slog.Logger
}

// NewLoggingDB creates a new LoggingDB wrapper
func NewLoggingDB(db *sql.DB, log *slog.Logger) *LoggingDB {
	return &LoggingDB{
		db:  db,
		log: log,
	}
}

// DB returns the underlying sql.DB
func (l *LoggingDB) DB() *sql.DB {
	return l.db
}

// QueryContext executes a query and logs the operation
func (l *LoggingDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	requestID := logger.GetRequestID(ctx)
	queryType := detectQueryType(query)

	rows, err := l.db.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	l.log.Info("query",
		slog.String("request_id", requestID),
		slog.String("type", queryType),
		slog.Duration("duration", duration),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		l.log.Error("query_error",
			slog.String("request_id", requestID),
			slog.String("type", queryType),
			slog.String("error", err.Error()),
			slog.Duration("duration", duration),
		)
	}

	return rows, err
}

// QueryRowContext executes a query that returns a single row and logs the operation
func (l *LoggingDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	start := time.Now()
	requestID := logger.GetRequestID(ctx)
	queryType := detectQueryType(query)

	row := l.db.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	l.log.Info("query",
		slog.String("request_id", requestID),
		slog.String("type", queryType),
		slog.Duration("duration", duration),
		slog.Bool("success", true),
	)

	return row
}

// ExecContext executes a query that doesn't return rows and logs the operation
func (l *LoggingDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	requestID := logger.GetRequestID(ctx)
	queryType := detectQueryType(query)

	result, err := l.db.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	var rowsAffected int64
	if result != nil {
		rowsAffected, _ = result.RowsAffected()
	}

	l.log.Info("exec",
		slog.String("request_id", requestID),
		slog.String("type", queryType),
		slog.Duration("duration", duration),
		slog.Int64("rows_affected", rowsAffected),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		l.log.Error("exec_error",
			slog.String("request_id", requestID),
			slog.String("type", queryType),
			slog.String("error", err.Error()),
			slog.Duration("duration", duration),
		)
	}

	return result, err
}

// BeginTx starts a transaction and logs the operation
func (l *LoggingDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*LoggingTx, error) {
	start := time.Now()
	requestID := logger.GetRequestID(ctx)

	tx, err := l.db.BeginTx(ctx, opts)
	duration := time.Since(start)

	l.log.Debug("tx_begin",
		slog.String("request_id", requestID),
		slog.Duration("duration", duration),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		l.log.Error("tx_begin_error",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return &LoggingTx{tx: tx, log: l.log, requestID: requestID, startTime: time.Now()}, nil
}

// Ping verifies the database connection
func (l *LoggingDB) Ping() error {
	return l.db.Ping()
}

// PingContext verifies the database connection with context
func (l *LoggingDB) PingContext(ctx context.Context) error {
	return l.db.PingContext(ctx)
}

// Close closes the database connection
func (l *LoggingDB) Close() error {
	return l.db.Close()
}

// LoggingTx wraps sql.Tx with logging
type LoggingTx struct {
	tx        *sql.Tx
	log       *slog.Logger
	requestID string
	startTime time.Time
}

// QueryContext executes a query within the transaction
func (t *LoggingTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	start := time.Now()
	queryType := detectQueryType(query)

	rows, err := t.tx.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	t.log.Info("tx_query",
		slog.String("request_id", t.requestID),
		slog.String("type", queryType),
		slog.Duration("duration", duration),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		t.log.Error("tx_query_error",
			slog.String("request_id", t.requestID),
			slog.String("type", queryType),
			slog.String("error", err.Error()),
		)
	}

	return rows, err
}

// QueryRowContext executes a query that returns a single row within the transaction
func (t *LoggingTx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	start := time.Now()
	queryType := detectQueryType(query)

	row := t.tx.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	t.log.Info("tx_query",
		slog.String("request_id", t.requestID),
		slog.String("type", queryType),
		slog.Duration("duration", duration),
		slog.Bool("success", true),
	)

	return row
}

// ExecContext executes a query that doesn't return rows within the transaction
func (t *LoggingTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	start := time.Now()
	queryType := detectQueryType(query)

	result, err := t.tx.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	var rowsAffected int64
	if result != nil {
		rowsAffected, _ = result.RowsAffected()
	}

	t.log.Info("tx_exec",
		slog.String("request_id", t.requestID),
		slog.String("type", queryType),
		slog.Duration("duration", duration),
		slog.Int64("rows_affected", rowsAffected),
		slog.Bool("success", err == nil),
	)

	if err != nil {
		t.log.Error("tx_exec_error",
			slog.String("request_id", t.requestID),
			slog.String("type", queryType),
			slog.String("error", err.Error()),
		)
	}

	return result, err
}

// Commit commits the transaction
func (t *LoggingTx) Commit() error {
	err := t.tx.Commit()
	duration := time.Since(t.startTime)

	if err != nil {
		t.log.Error("tx_commit_error",
			slog.String("request_id", t.requestID),
			slog.String("error", err.Error()),
			slog.Duration("tx_duration", duration),
		)
	} else {
		t.log.Debug("tx_commit",
			slog.String("request_id", t.requestID),
			slog.Duration("tx_duration", duration),
		)
	}

	return err
}

// Rollback rolls back the transaction
func (t *LoggingTx) Rollback() error {
	err := t.tx.Rollback()
	duration := time.Since(t.startTime)

	// Don't log error for already committed transactions
	if err != nil && err != sql.ErrTxDone {
		t.log.Warn("tx_rollback",
			slog.String("request_id", t.requestID),
			slog.String("error", err.Error()),
			slog.Duration("tx_duration", duration),
		)
	} else if err == nil {
		t.log.Debug("tx_rollback",
			slog.String("request_id", t.requestID),
			slog.Duration("tx_duration", duration),
		)
	}

	return err
}

// detectQueryType determines the type of SQL query from the query string
func detectQueryType(query string) string {
	query = strings.TrimSpace(strings.ToUpper(query))

	switch {
	case strings.HasPrefix(query, "SELECT"):
		return "SELECT"
	case strings.HasPrefix(query, "INSERT"):
		return "INSERT"
	case strings.HasPrefix(query, "UPDATE"):
		return "UPDATE"
	case strings.HasPrefix(query, "DELETE"):
		return "DELETE"
	case strings.HasPrefix(query, "CREATE"):
		return "CREATE"
	case strings.HasPrefix(query, "ALTER"):
		return "ALTER"
	case strings.HasPrefix(query, "DROP"):
		return "DROP"
	default:
		return "OTHER"
	}
}
