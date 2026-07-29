package database

import (
	"log"

	"gojo/internal/config"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/v4"
)

func StartMigrations(cfg *config.Config) {
	migration, err := migrate.New("file://migrations", cfg.DatabaseURLMigration)
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.Up(); err != nil && migrate.ErrNoChange {
		log.Fatal(err)
	}
}
