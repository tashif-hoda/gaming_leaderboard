package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// MigrateDB handles database migrations
func (db *DB) MigrateDB(down bool, targetVersion int) error {
	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create the postgres driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+os.Getenv("MIGRATIONS_PATH"),
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("migration failed: %v", err)
	}

	// Get current version before migration
	currentVersion, _, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get current version: %v", err)
	}

	if targetVersion >= 0 {
		// Migrate to specific version
		if err := m.Migrate(uint(targetVersion)); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not migrate to version %d: %v", targetVersion, err)
		}
		log.Printf("Migrated from version %d to version %d", currentVersion, targetVersion)
	} else if down {
		// Down migrate all
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not run down migrations: %v", err)
		}
		log.Printf("Down migrated from version %d", currentVersion)
	} else {
		// Up migrate all
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("could not run up migrations: %v", err)
		}
		newVersion, _, _ := m.Version()
		log.Printf("Up migrated from version %d to version %d", currentVersion, newVersion)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get migration version: %v", err)
	}

	log.Printf("Current migration version: %d (dirty: %v)", version, dirty)
	return nil
}
