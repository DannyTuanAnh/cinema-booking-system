package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DSN string
}

// NewConnection creates a new PostgreSQL database connection
func NewConnection(config Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// DefaultConfig returns a default database configuration for local development
func DefaultConfig() Config {
	// Ưu tiên DATABASE_URL (Render / Production)
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return Config{
			DSN: dsn + "?sslmode=require",
		}
	}

	// Fallback cho local
	return Config{
		DSN: fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			"localhost",
			5432,
			"postgres",
			"1",
			"cinema",
		),
	}

}
