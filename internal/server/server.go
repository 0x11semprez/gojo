package server

import (
	"net/http"
	"time"

	"gojo/internal/config"
)

type Server struct {
	*http.Server
}

func NewServe(c *config.Config, mux http.Handler) (*Server, error) {
	return &Server{
		&http.Server{
			Addr:              c.Port,
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20,
		},
	}, nil
}
