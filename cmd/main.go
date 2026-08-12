// Package main est le point d'entrée de l'application gojo.
// Il assemble la configuration, la base de données et le serveur HTTP,
// puis démarre l'écoute des requêtes entrantes.
package main

import (
	"log"
	"net/http"

	"gojo/internal/config"
	"gojo/internal/database"
	"gojo/internal/server"
)

func main() {
	// Charge la configuration (variables d'environnement via .env).
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// mux est le routeur HTTP standard qui associera les chemins d'URL
	// aux handlers (voir internal/user pour les handlers métier).
	mux := http.NewServeMux()
	srv, err := server.NewServe(cfg, mux)
	if err != nil {
		log.Fatal(err)
	}

	// Applique les migrations de base de données au démarrage,
	// pour garantir que le schéma est à jour avant de servir des requêtes.
	database.StartMigrations(cfg)

	// Démarre le serveur HTTP en mode bloquant.
	srv.ListenAndServe()
}
