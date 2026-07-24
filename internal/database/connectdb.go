package database

import (
	"database/sql"
	"log"

	"gojo/internal/config"
)

func ConnectDB(databaseURL *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db, err
}
