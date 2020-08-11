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

func setupRoutes(handler *handler.Handler) {
	http.HandleFunc("/signup", handleError(handler.SignUp))
	http.HandleFunc("/login", handleError(handler.Login))
	http.HandleFunc("/refresh", handleError(handler.RefreshToken))

	http.HandleFunc("/welcome", handleError(withAuth(handler.Welcome)))

	http.HandleFunc("/create-list-item", handleError(withAuth(handler.CreateListItem)))
	http.HandleFunc("/list-items-for-user", handleError(withAuth(handler.ReadListItemsForUser)))
}

func getConnectionInfo() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "undefined"
	}
	fmt.Println("🌍 env:", env)

	switch env {
	case "production":
		return os.Getenv("DATABASE_URL")
	default:
		const (
			host     = "localhost"
			port     = 5432
			user     = "naz"
			password = "test"
			dbname   = "filmlistclub"
			sslmode  = "disable"
		)

		return fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("🗼 Listening on http://localhost:" + port)

	return ":" + port
}

// --- middleware ---
// TODO: best pattern for moving to another file? middleware package?

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
