package database

import (
	"database/sql"
	"log"

	"gojo/internal/app"
)

func connectDB(databaseURL *app.App) {
	db, err := sql.Open("postgres", databaseURL.Config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
