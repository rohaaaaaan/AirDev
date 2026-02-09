package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %v", err)
	}

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %v", err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}
