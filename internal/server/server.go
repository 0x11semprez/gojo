package server

import (
	"log"
	"net/http"
	"time"

	"gojo/internal/app"
)

type Server struct {
	*http.Server
}

func NewServe(c *app.Config) Server {
	return Server{
		&http.Server{
			Addr:              c.Port,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}
}

func StartServer(s *Server) {
	err := http.ListenAndServe(s.Addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
