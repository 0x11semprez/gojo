// Package database gère la connexion à PostgreSQL (via l'ORM Bun) ainsi
// que l'exécution des migrations de schéma au démarrage de l'application.
package database

import (
	"database/sql"
	"errors"
	"time"

	"gojo/internal/config"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// ConnectDB crée et configure une connexion PostgreSQL en utilisant l'ORM Bun.
//
// Elle retourne :
//   - *bun.DB : l'instance Bun utilisée pour exécuter les requêtes
//   - error : une erreur si la configuration ou la connexion échoue
func ConnectDB(dsn *config.Config) (*bun.DB, error) {
	// Vérifie que la configuration existe.
	// Un pointeur nil signifie qu'aucune configuration n'a été fournie.
	if dsn == nil {
		return nil, errors.New("app is nil")
	}

	// Vérifie que l'URL de la base de données est renseignée.
	// databaseURL contient des informations comme :
	// postgres://user:password@host:port/database
	if dsn.DatabaseURL == "" {
		return nil, errors.New("database URL is empty")
	}

	// Crée une connexion SQL standard en utilisant le driver PostgreSQL de Bun.
	// Cela n'ouvre pas immédiatement de connexion réseau ;
	// les connexions sont créées à la demande, lors de l'exécution des requêtes.
	pgdb := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(dsn.DatabaseURL),
		),
	)

	// Configure le pool de connexions à la base de données.
	//
	// Nombre maximum de connexions pouvant être ouvertes simultanément.
	pgdb.SetMaxOpenConns(25)

	// Nombre maximum de connexions inactives (inutilisées mais prêtes) gardées dans le pool.
	pgdb.SetMaxIdleConns(10)

	// Ferme les connexions inactives qui n'ont pas été utilisées depuis 5 minutes.
	pgdb.SetConnMaxIdleTime(5 * time.Minute)

	// Recycle les connexions après 25 minutes.
	// Cela évite de conserver indéfiniment des connexions anciennes ou instables.
	pgdb.SetConnMaxLifetime(25 * time.Minute)

	// Crée l'instance de l'ORM Bun.
	// Bun utilisera le dialecte PostgreSQL pour générer les requêtes SQL.
	db := bun.NewDB(
		pgdb,
		pgdialect.New(),
	)

	// Retourne l'instance de base de données configurée.
	return db, nil
}
