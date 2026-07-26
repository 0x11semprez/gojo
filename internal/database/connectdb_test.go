package database

import (
	"log"
	"testing"

	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/server"
)

func TestConnectDb(t *testing.T) {
	config := config.NewConfig()
	s := server.NewServe(config)
	db, err := ConnectDB(config)
	if err != nil {
		log.Fatal(err)
	}
	app := app.NewApp(&s, db, config)
}
