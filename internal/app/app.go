package app

import (
	"gojo/internal/config"
	"gojo/internal/server"

	"github.com/uptrace/bun"
)

type App struct {
	Server   *server.Server
	Database *bun.DB
	Config   *config.Config
}

func NewApp(s *server.Server, db *bun.DB, c *config.Config) *App {
	return &App{
		Server:   s,
		Database: db,
		Config:   c,
	}
}
