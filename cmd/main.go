package main

import (
	"fmt"
	"log"

	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/database"
	"gojo/internal/server"
)

func main() {
	config := config.NewConfig()

	s := server.NewServe(config)

	server.Run()

	server.StartServer(&s)

	db, err := database.ConnectDB(config)
	if err != nil {
		log.Fatal(err)
	}

	application := app.NewApp(&s, db, config)

	fmt.Println(application.Config.Port)
}
