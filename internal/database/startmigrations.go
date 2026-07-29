package database

import (
	"errors"
	"log"

	"gojo/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func StartMigrations(cfg *config.Config) {
	migration, err := migrate.New("file://migrations", cfg.DatabaseURLMigration)
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}
