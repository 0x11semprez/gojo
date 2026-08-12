package database

import (
	"errors"
	"log"
	"path/filepath"
	"runtime"

	"gojo/internal/config"

	"github.com/golang-migrate/migrate/v4"
	// Anonymous imports (side effects only): these two packages
	// register their drivers with golang-migrate (destination
	// "postgres" and source "file"), without us directly calling
	// an exported function from the package.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// StartMigrations applies the pending database migrations
// (SQL files from the internal/database/migrations folder) using
// golang-migrate. It is called at application startup
// (see cmd/main.go) to guarantee the schema is up to date.
// On failure, it stops the process (log.Fatal), since starting
// with an out-of-sync schema would be dangerous.
func StartMigrations(cfg *config.Config) {
	// "file://migrations" is resolved relative to the process's
	// working directory, which depends on where `go run`/the binary
	// is launched from and therefore does not reliably point to
	// internal/database/migrations. We anchor the path to this
	// source file so migrations are found regardless of the current
	// working directory.
	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "migrations")

	migration, err := migrate.New("file://"+migrationsDir, cfg.DatabaseURLMigration)
	if err != nil {
		log.Fatalf("failed to initialize database migrations: %v", err)
	}

	// migrate.ErrNoChange is not a real error: it simply means there
	// was no pending migration to apply.
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to apply database migrations: %v", err)
	}

	log.Println("database migrations applied successfully")
}
