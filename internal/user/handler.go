package user

import (
	"fmt"
	"log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "All is good")
	if err != nil {
		log.Fatal(err)
	}
}
