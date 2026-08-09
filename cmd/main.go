package main

import (
	"log"
	"net/http"

	"gojo/internal/config"
	"gojo/internal/database"
	"gojo/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	srv, err := server.NewServe(cfg, mux)
	if err != nil {
		log.Fatal(err)
	}

	database.StartMigrations(cfg)

	srv.ListenAndServe()
}
