package database

import (
	"database/sql"
	"errors"
	"time"

	"gojo/internal/app"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// ConnectDB creates and configures a PostgreSQL database connection using Bun ORM.
//
// It returns:
//   - *bun.DB: the Bun database instance used for queries
//   - error: an error if the configuration or connection setup fails
func ConnectDB(dsn *app.App) (*bun.DB, error) {
	// Check if the application instance exists.
	// A nil pointer means no application configuration was provided.
	if dsn == nil {
		return nil, errors.New("app is nil")
	}

	// Check if the configuration object exists.
	// Without it, we cannot access the database URL.
	if dsn.Config == nil {
		return nil, errors.New("app config is nil")
	}

	// Check if the database connection string is configured.
	// The DSN contains information such as:
	// postgres://user:password@host:port/database
	if dsn.Config.DatabaseURL == "" {
		return nil, errors.New("database URL is empty")
	}

	// Create a standard SQL database connection using Bun's PostgreSQL driver.
	// This does not immediately open a connection;
	// connections are created lazily when queries are executed.
	pgdb := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(dsn.Config.DatabaseURL),
		),
	)

	// Configure the database connection pool.
	//
	// Maximum number of connections that can be opened at the same time.
	pgdb.SetMaxOpenConns(25)

	// Maximum number of idle (unused but ready) connections kept in the pool.
	pgdb.SetMaxIdleConns(10)

	// Close idle connections that have not been used for 5 minutes.
	pgdb.SetConnMaxIdleTime(5 * time.Minute)

	// Recycle connections after 25 minutes.
	// This prevents keeping old or unhealthy connections forever.
	pgdb.SetConnMaxLifetime(25 * time.Minute)

	// Create the Bun ORM instance.
	// Bun will use PostgreSQL dialect to generate SQL queries.
	db := bun.NewDB(
		pgdb,
		pgdialect.New(),
	)

	// Return the configured database instance.
	return db, nil
}
