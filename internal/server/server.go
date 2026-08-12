// Package server encapsule la configuration du serveur HTTP standard de
// l'application (adresse d'écoute, délais d'expiration, routeur).
package server

import (
	"net/http"
	"time"

	"gojo/internal/config"
)

// Server enrobe http.Server. L'embarquement (composition) permet
// d'appeler directement les méthodes de http.Server (ex: ListenAndServe)
// sur une valeur de type Server, tout en gardant la possibilité
// d'ajouter des champs/méthodes propres à l'application plus tard.
type Server struct {
	*http.Server
}

// NewServe construit un Server prêt à démarrer, avec des délais
// d'expiration (timeouts) définis pour éviter que des connexions lentes
// ou inactives ne bloquent le serveur indéfiniment.
func NewServe(c *config.Config, mux http.Handler) (*Server, error) {
	return &Server{
		&http.Server{
			Addr:    c.Port,
			Handler: mux,
			// Délai maximum pour lire l'intégralité de la requête (corps inclus).
			ReadTimeout: 5 * time.Second,
			// Délai maximum pour lire uniquement les en-têtes de la requête.
			ReadHeaderTimeout: 2 * time.Second,
			// Délai maximum pour écrire la réponse.
			WriteTimeout: 5 * time.Second,
			// Délai avant fermeture d'une connexion keep-alive inactive.
			IdleTimeout: 60 * time.Second,
			// Taille maximale des en-têtes de requête acceptés (1 Mo).
			MaxHeaderBytes: 1 << 20,
		},
	}, nil
}
