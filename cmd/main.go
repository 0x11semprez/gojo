package main

import (
	"gojo/internal/app"
	"gojo/internal/server"
)

func main() {
	config := app.NewConfig()

	s := server.NewServe(config)

	server.Run()

	server.StartServer(&s)
}
