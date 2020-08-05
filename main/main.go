package main

import (
	"log"
	"net/http"

	"github.com/FilmListClub/backend/api"
)

func main() {
	// "Login" and "Welcome" are the handlers that we will implement
	http.HandleFunc("/login", api.Login)
	http.HandleFunc("/welcome", api.Welcome)
	http.HandleFunc("/refresh", api.Refresh)

	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}
