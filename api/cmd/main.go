// Package main is the entry point of the gojo application.
// It wires together the configuration, the database and the HTTP server,
// then starts listening for incoming requests.
package main

import (
	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/database"
	"gojo/internal/server"
	"log"
	"net/http"
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

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// a is the dependency container injected into every HTTP handler
	// (see cmd/rooter.go).
	a := app.App{}.NewApp(srv, db, cfg)
	registerRoutes(mux, a)

	log.Printf("server is running and listening on %s", cfg.Port)

	// Start the HTTP server in blocking mode.
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped unexpectedly: %v", err)
	}
}
