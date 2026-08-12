package database

import (
	"errors"
	"log"
	"path/filepath"
	"runtime"

	"gojo/internal/config"

	"github.com/golang-migrate/migrate/v4"
	// Imports anonymes (effets de bord uniquement) : ces deux packages
	// enregistrent leurs drivers auprès de golang-migrate (destination
	// "postgres" et source "file"), sans qu'on appelle directement
	// une fonction exportée du package.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// StartMigrations applique les migrations de base de données en attente
// (fichiers SQL du dossier internal/database/migrations) à l'aide de
// golang-migrate. Elle est appelée au démarrage de l'application
// (voir cmd/main.go) pour garantir que le schéma est à jour.
// En cas d'échec, elle arrête le processus (log.Fatal), car démarrer
// avec un schéma désynchronisé serait dangereux.
func StartMigrations(cfg *config.Config) {
	// "file://migrations" est résolu par rapport au répertoire de travail du
	// processus, qui dépend de l'endroit où `go run`/le binaire est lancé et
	// ne pointe donc pas de façon fiable vers internal/database/migrations.
	// On ancre le chemin à ce fichier source pour que les migrations soient
	// trouvées quel que soit le répertoire courant.
	_, thisFile, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(thisFile), "migrations")

	migration, err := migrate.New("file://"+migrationsDir, cfg.DatabaseURLMigration)
	if err != nil {
		log.Fatal(err)
	}

	// migrate.ErrNoChange n'est pas une vraie erreur : elle signifie
	// simplement qu'il n'y avait aucune migration en attente à appliquer.
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}
