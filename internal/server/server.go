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

func NewServe(c *app.Config, handler http.Handler) Server {
	return Server{
		&http.Server{
			Addr:              c.Port,
			Handler:           handler,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}
}

func StartServer(s *Server) {
	err := http.ListenAndServe(s.Addr, s.Handler)
	if err != nil {
		log.Fatal(err)
	}
}
