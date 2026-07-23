package app

import (
	"database/sql"
	"net/http"
)

type App struct {
	Server   *http.Server
	Database *sql.DB
	Config   *Config
}

func NewApp(s *http.Server, db *sql.DB, c *Config) *App {
	return &App{
		Server:   s,
		Database: db,
		Config:   c,
	}
}
