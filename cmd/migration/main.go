package migration

import (
	"embed"
	"log"

	"github.com/uptrace/bun/migrate"
)

var (
	sqlMigrations embed.FS
	Migrations    = migrate.NewMigrations()
)

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		log.Fatal(err)
	}
}
