package server

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	*http.Server
}

func NewServe(c *Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              c.Port,
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

func (s *Server) StartServer() {
	err := http.ListenAndServe(s.Addr, s.Handler)
	if err != nil {
		fmt.Errorf("There is a problem here: %w", err)
	}
}
