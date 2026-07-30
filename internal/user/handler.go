package user

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"gojo/internal/app"
	"gojo/internal/middleware"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "All is good")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(db *app.App) (string, error) {
	secret, err := middleware.GenerateSecret()
	if err != nil {
		errors.New("cannot create user")
	}

	fmt.Print("your secret is: %v", secret)

	return secret, nil
}
