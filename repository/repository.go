// package repository implements a database connection and repositories for each table.

package repository

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
	"lopa.to/sonimulus/config"
)

// NewDB creates a new PostgreSQL database connection.
func NewDB(config *config.Config) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", config.DBUrl)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return nil, err
	}
	return db, nil
}
