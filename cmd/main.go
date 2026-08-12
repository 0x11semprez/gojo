// Package main is the entry point of the gojo application.
// It wires together the configuration, the database and the HTTP server,
// then starts listening for incoming requests.
package main

import (
	"log"
	"net/http"

	"gojo/internal/config"
	"gojo/internal/database"
	"gojo/internal/server"
)

func main() {
	// Load the configuration (environment variables via .env).
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// mux is the standard HTTP router that maps URL paths
	// to handlers (see internal/user for the business handlers).
	mux := http.NewServeMux()
	srv, err := server.NewServe(cfg, mux)
	if err != nil {
		log.Fatalf("failed to initialize HTTP server: %v", err)
	}

	// Apply database migrations at startup, to guarantee the schema
	// is up to date before serving any request.
	database.StartMigrations(cfg)

	log.Printf("server is running and listening on %s", cfg.Port)

	// Start the HTTP server in blocking mode.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped unexpectedly: %v", err)
	}
}
