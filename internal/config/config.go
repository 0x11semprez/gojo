// Package config charge la configuration de l'application depuis des
// fichiers .env et l'expose sous forme de structure typée (Config).
package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config regroupe les paramètres de configuration nécessaires au
// démarrage de l'application (accès base de données, port d'écoute).
type Config struct {
	// DatabaseURL est le DSN utilisé pour les requêtes applicatives normales.
	DatabaseURL string
	// DatabaseURLMigration est le DSN utilisé spécifiquement pour exécuter
	// les migrations (peut nécessiter des droits différents de DatabaseURL).
	DatabaseURLMigration string
	// Port est le port sur lequel le serveur HTTP écoute.
	Port string
}

// NewConfig charge le fichier .env à la racine du projet et construit
// la configuration de production à partir des variables d'environnement.
func NewConfig() (*Config, error) {
	// "../../.env" est résolu relativement au répertoire de travail du
	// processus au moment de l'exécution, pas au dossier de ce fichier.
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, errors.New("cannot load	.env file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION"),
		Port:                 os.Getenv("PORT"),
	}, nil
}

// NewConfigTest charge le fichier .env.test et construit une configuration
// dédiée aux tests, isolée de la configuration de production (base de
// données de test, DSN de migration de test).
func NewConfigTest() (*Config, error) {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		return nil, errors.New("cannot load .env.test file")
	}
	return &Config{
		DatabaseURL:          os.Getenv("DATABASE_URL_TEST"),
		DatabaseURLMigration: os.Getenv("DATABASE_URL_MIGRATION_TEST"),
		Port:                 os.Getenv("PORT"),
	}, nil
}
