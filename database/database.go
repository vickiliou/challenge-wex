package database

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"golang.org/x/exp/slog"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Setup initializes and configures the SQLite database.
func Setup() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "wex.db")
	if err != nil {
		slog.Error("Failed to open SQLite database", "error", err.Error())
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		slog.Error("Failed to run database migrations", "error", err.Error())
		return nil, err
	}

	return db, nil
}

// runMigrations runs database migrations using the provided database connection and configuration.
func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		slog.Error("Failed to set SQLite dialect", "error", err.Error())
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		slog.Error("Failed to run migrations", "error", err.Error())
		return err
	}

	return nil
}
