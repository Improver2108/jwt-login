package db

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
)

func RunMigration(databaseURL string) error {
	m, err := migrate.New("file://db/migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator : %w", err)
	}
	defer m.Close()

	beforeVersion, _, _ := m.Version()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("Migration: already up to date, nothing to run")
			return nil
		}
		log.Printf("Migration failed: %v", err)
		log.Printf("Rolling back to version %d...", beforeVersion)

		if beforeVersion == 0 {
			if downErr := m.Down(); downErr != nil && downErr != migrate.ErrNoChange {
				return fmt.Errorf("migration failed AND rollback failed: %w", downErr)
			}
		} else {
			if migrateErr := m.Migrate(beforeVersion); migrateErr != nil && migrateErr != migrate.ErrNoChange {
				return fmt.Errorf("migration failed AND rollback failed: %w", migrateErr)
			}
		}
		return fmt.Errorf("migrations rolled back to version %d due to error: %w", beforeVersion, err)
	}
	afterVersion, _, _ := m.Version()
	log.Printf("Migrations applied successfully (version %d -> %d)", beforeVersion, afterVersion)
	return nil
}
