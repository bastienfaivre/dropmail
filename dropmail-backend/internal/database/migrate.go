package database

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all pending database migrations.
// migrationsPath should be a file:// URL to the migrations directory.
func (db *DB) RunMigrations(migrationsPath string) error {
	driver, err := sqlite3.WithInstance(db.conn, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Get current version before migration
	version, dirty, _ := m.Version()
	slog.Info("current migration version", "version", version, "dirty", dirty)

	// Run all pending migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Get version after migration
	newVersion, _, _ := m.Version()
	if newVersion != version {
		slog.Info("migrations applied successfully", "from_version", version, "to_version", newVersion)
	} else {
		slog.Info("database schema is up to date", "version", newVersion)
	}

	return nil
}
