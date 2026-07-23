package server

import (
	"net/http"

	"gojo/internal/user"
)

var mux = http.NewServeMux()

func Run() {
	mux.HandleFunc("GET /health", user.Health)
}
