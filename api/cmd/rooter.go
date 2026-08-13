// Package main (rooter.go file) declares the HTTP routing
// configuration of the application: it maps each route to its
// handler, using the method+pattern syntax supported by the standard
// library's http.ServeMux (Go 1.22+).
package main

import (
	"gojo/internal/app"
	"gojo/internal/user"
	"gojo/internal/wallet"
	"net/http"
)

// registerRoutes wires every HTTP route the application exposes onto
// mux, binding each handler to the shared App (server, database,
// config).
func registerRoutes(mux *http.ServeMux, a *app.App) {
	mux.HandleFunc("GET /health", user.Health)

	mux.HandleFunc("POST /users", user.CreateUserHandler(a))
	mux.HandleFunc("PATCH /users/{id}", user.UpdateUserHandler(a))
	mux.HandleFunc("DELETE /users/{id}", user.DeleteUserHandler(a))

	// Verifies a username/password pair and returns the matching
	// user's id, required before generating a wallet (see POST
	// /wallets below).
	mux.HandleFunc("POST /login", user.LoginHandler(a))

	// Generates a fresh keypair for the given network ("bitcoin" or
	// "monero" in the request body) and stores it envelope-encrypted.
	mux.HandleFunc("POST /wallets", wallet.CreateWalletHandler(a))
}
