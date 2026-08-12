// Package app groups the core dependencies of the application
// (HTTP server, database connection, configuration) into a single
// struct, which is then injected into the business handlers (e.g. the user package).
package app

import (
	"gojo/internal/config"
	"gojo/internal/server"

	"github.com/uptrace/bun"
)

// App is the application's dependency container. It lets handlers
// receive a single access point to the server, the database and the
// config, instead of multiplying parameters.
type App struct {
	Server   *server.Server
	Database *bun.DB
	Config   *config.Config
}

// NewApp builds a new App instance from already initialized
// components (server, database, configuration).
func (a App) NewApp(s *server.Server, db *bun.DB, c *config.Config) *App {
	return &App{
		Server:   s,
		Database: db,
		Config:   c,
	}
}
