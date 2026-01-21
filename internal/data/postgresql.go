package data

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new PostgreSQL database connection.
func NewPostgresDB(uri string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", uri)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, err
	}

	slog.Info("Database connection established")

	return db, nil
}
