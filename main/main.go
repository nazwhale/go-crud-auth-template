package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/FilmListClub/backend/auth"

	"github.com/FilmListClub/backend/dao"
	"github.com/FilmListClub/backend/handler"
	_ "github.com/lib/pq"
)

// TODO: extract credentials stuff into middleware
// authenticate(handler.SignUp) ???
func main() {
	db, err := sql.Open("postgres", getConnectionInfo())
	if err != nil {
		log.Fatalf("db crashed: %q", err)
	}
	defer db.Close()

	dao := dao.New(db)
	handler := handler.New(dao)

	setupRoutes(handler)
	log.Fatal(http.ListenAndServe(getPort(), nil))
}

// --- db connection ---

const (
	host     = "localhost"
	port     = 5432
	user     = "naz" // or "postgres"?
	password = "test"
	dbname   = "filmlistclub"
	sslmode  = "disable"
)

func getConnectionInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Listening on http://localhost:" + port)

	return ":" + port
}

// --- routes & middleware ---

func setupRoutes(handler *handler.Handler) {
	http.HandleFunc("/signup", handleError(handler.SignUp))
	http.HandleFunc("/login", handleError(handler.Login))
	http.HandleFunc("/refresh", handleError(handler.RefreshToken))

	http.HandleFunc("/welcome", handleError(withAuth(handler.Welcome)))
}

type H func(http.ResponseWriter, *http.Request) *handler.Error
type httpFunc func(http.ResponseWriter, *http.Request)

func handleError(h H) httpFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := h(w, r)
		if e != nil {
			log.Printf("HTTP %d - %s", e.ResponseCode, e.Message)
			http.Error(w, e.Message, e.ResponseCode)
		}
	}
}

func withAuth(h H) H {
	return func(w http.ResponseWriter, r *http.Request) *handler.Error {
		_, err := auth.ValidateRequest(r)
		if err != nil {
			e := &handler.Error{
				Message:      err.Error(),
				ResponseCode: http.StatusUnauthorized,
			}
			return e.Wrap("error validating request")
		}

		return h(w, r)
	}
}
