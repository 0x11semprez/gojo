package database

import (
	"database/sql"

	"gojo/internal/app"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func ConnectDB(dsn *app.App) *bun.DB {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn.Config.DatabaseURL)))

	return bun.NewDB(pgdb, pgdialect.New())
}
