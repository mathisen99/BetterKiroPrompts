package db

import (
	"context"
	"database/sql"
	"log"
	"os"

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

	if err := DB.PingContext(ctx); err != nil {
		return err
	}

	log.Println("Database connected successfully")
	return nil
}
