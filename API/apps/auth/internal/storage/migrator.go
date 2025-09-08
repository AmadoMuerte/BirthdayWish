package storage

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func ApplyMigrations(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migrate driver: %w", err)
	}

	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	migrationPath := filepath.Join(basepath, "migrations")

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("could not get migration version: %w", err)
	}

	if dirty {
		log.Printf("Database is dirty at version %d. Attempting to fix...", version)
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("could not force version %d: %w", version, err)
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not apply migrations: %w", err)
	}

	currentVersion, _, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("could not get current migration version: %w", err)
	}

	log.Printf("Database migrations are up to date. Current version: %d", currentVersion)
	return nil
}
