package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gojo/internal/app"
	"gojo/internal/middleware"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Print(w, "All is good\n")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateUser(db *app.App) (string, error) {
	secret, err := middleware.GenerateSecret()
	if err != nil {
		errors.New("cannot create user")
	}

	ctx := context.Background()

	user := &User{Secret: secret}
	_, err = db.Database.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("your secret is: %v", secret)

	return secret, nil
}

func DeleteUser(db *app.App, secret string) (string, error) {
	var selectedUser User

	ctx := context.Background()

	_, err := db.Database.NewDelete().Model(selectedUser).Where("secret = ?", secret).Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	sucess := "user is delete"
	return sucess, nil
}
