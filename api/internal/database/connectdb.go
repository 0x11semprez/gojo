// Package database handles the connection to PostgreSQL (via the Bun
// ORM) as well as running schema migrations at application startup.
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

// ConnectDB creates and configures a PostgreSQL connection using the Bun ORM.
//
// It returns:
//   - *bun.DB: the Bun instance used to run queries
//   - error: an error if the configuration or the connection fails
func ConnectDB(dsn *config.Config) (*bun.DB, error) {
	// Check that the configuration exists.
	// A nil pointer means no configuration was provided.
	if dsn == nil {
		return nil, errors.New("app is nil")
	}

	// Check that the database URL is set.
	// databaseURL contains information like:
	// postgres://user:password@host:port/database
	if dsn.DatabaseURL == "" {
		return nil, errors.New("database URL is empty")
	}

	// Create a standard SQL connection using Bun's PostgreSQL driver.
	// This does not immediately open a network connection;
	// connections are created on demand, when queries are executed.
	pgdb := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(dsn.DatabaseURL),
		),
	)

	// Configure the database connection pool.
	//
	// Maximum number of connections that can be open simultaneously.
	pgdb.SetMaxOpenConns(25)

	// Maximum number of idle (unused but ready) connections kept in the pool.
	pgdb.SetMaxIdleConns(10)

	// Close idle connections that have not been used for 5 minutes.
	pgdb.SetConnMaxIdleTime(5 * time.Minute)

	// Recycle connections after 25 minutes.
	// This avoids keeping old or potentially unstable connections indefinitely.
	pgdb.SetConnMaxLifetime(25 * time.Minute)

	// Create the Bun ORM instance.
	// Bun will use the PostgreSQL dialect to generate SQL queries.
	db := bun.NewDB(
		pgdb,
		pgdialect.New(),
	)

	// Return the configured database instance.
	return db, nil
}
