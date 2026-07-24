package app

import (
	"database/sql"

	"gojo/internal/config"
	"gojo/internal/server"
)

type App struct {
	Server   *server.Server
	Database *sql.DB
	Config   *config.Config
}

func NewApp(s *server.Server, db *sql.DB, c *config.Config) *App {
	return &App{
		Server:   s,
		Database: db,
		Config:   c,
	}
}
