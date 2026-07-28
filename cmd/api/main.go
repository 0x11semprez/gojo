package main

import (
	"log"

	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/database"
	"gojo/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServe(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	application := app.NewApp(s, db, cfg)

	server.Run()

	server.StartServer(s)
}
