// Package server encapsulates the standard HTTP server configuration
// of the application (listen address, timeouts, router).
package server

import (
	"net/http"
	"time"

	"gojo/internal/config"
)

// Server wraps http.Server. Embedding (composition) allows calling
// http.Server methods directly (e.g. ListenAndServe) on a Server
// value, while keeping the ability to add application-specific
// fields/methods later.
type Server struct {
	*http.Server
}

// NewServe builds a Server ready to start, with timeouts configured
// to prevent slow or idle connections from blocking the server indefinitely.
func NewServe(c *config.Config, mux http.Handler) (*Server, error) {
	return &Server{
		&http.Server{
			Addr:    c.Port,
			Handler: mux,
			// Maximum duration to read the entire request (body included).
			ReadTimeout: 5 * time.Second,
			// Maximum duration to read only the request headers.
			ReadHeaderTimeout: 2 * time.Second,
			// Maximum duration to write the response.
			WriteTimeout: 5 * time.Second,
			// Duration before closing an idle keep-alive connection.
			IdleTimeout: 60 * time.Second,
			// Maximum size of accepted request headers (1 MB).
			MaxHeaderBytes: 1 << 20,
		},
	}, nil
}
