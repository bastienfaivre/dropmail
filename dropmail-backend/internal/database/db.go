package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the database connection.
type DB struct {
	conn *sql.DB
	path string
}

// Open opens a connection to the SQLite database at the given path.
// Creates the database file and parent directories if they don't exist.
func Open(dbPath string) (*DB, error) {
	// Ensure parent directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open SQLite connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Verify connection
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable WAL mode for better concurrent access
	if _, err := conn.Exec("PRAGMA journal_mode=WAL"); err != nil {
		slog.Warn("failed to enable WAL mode", "error", err)
	}

	// Enable foreign keys
	if _, err := conn.Exec("PRAGMA foreign_keys=ON"); err != nil {
		slog.Warn("failed to enable foreign keys", "error", err)
	}

	slog.Info("database connection established", "path", dbPath)

	return &DB{conn: conn, path: dbPath}, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	if db.conn != nil {
		slog.Info("closing database connection")
		return db.conn.Close()
	}
	return nil
}

// Conn returns the underlying sql.DB connection.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// Path returns the database file path.
func (db *DB) Path() string {
	return db.path
}
