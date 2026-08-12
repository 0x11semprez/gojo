package database

import (
	"errors"
	"log"
	"path/filepath"
	"runtime"

	"gojo/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func StartMigrations(cfg *config.Config) {
	// "file://migrations" is resolved against the process working directory,
	// which depends on where `go run`/the binary is launched from and does
	// not reliably point at internal/database/migrations. Anchor the path to
	// this source file instead so migrations are found regardless of cwd.
	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "migrations")

	migration, err := migrate.New("file://"+migrationsDir, cfg.DatabaseURLMigration)
	if err != nil {
		log.Fatal(err)
	}

	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}
