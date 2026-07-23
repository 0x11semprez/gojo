package app

import "net/http"

type App struct {
	Server   *http.Server
	Database *sql.db
	Config   Config
}
