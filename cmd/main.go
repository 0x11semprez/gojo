package main

import (
	"net/http"

	"gojo/internal/app"
	"gojo/internal/server"
)

func main() {
	config := app.NewConfig()

	mux := http.NewServeMux()

	s := server.NewServe(config, mux)

	server.StartServer(&s)
}
