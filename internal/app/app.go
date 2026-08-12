// Package app regroupe les dépendances centrales de l'application
// (serveur HTTP, connexion base de données, configuration) dans une
// seule structure, injectée ensuite aux handlers métier (ex: package user).
package app

import (
	"gojo/internal/config"
	"gojo/internal/server"

	"github.com/uptrace/bun"
)

// App est le conteneur de dépendances de l'application. Il permet de
// passer aux handlers un accès unique au serveur, à la base et à la config,
// plutôt que de multiplier les paramètres.
type App struct {
	Server   *server.Server
	Database *bun.DB
	Config   *config.Config
}

// NewApp construit une nouvelle instance d'App à partir des composants
// déjà initialisés (serveur, base de données, configuration).
func (a App) NewApp(s *server.Server, db *bun.DB, c *config.Config) *App {
	return &App{
		Server:   s,
		Database: db,
		Config:   c,
	}
}
