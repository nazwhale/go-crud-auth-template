package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FilmListClub/backend/dao"
	"golang.org/x/crypto/bcrypt"
)

// Create a struct to read the username and password from the request body
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) *Error {
	// TODO: pull CORS out into middleware
	(w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	credentials := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(credentials)
	if err != nil {
		return &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusBadRequest,
		}
	}

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8
	// (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 8)
	if err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
		return e.Wrap("error hashing password")
	}

	// Create a user
	user := dao.User{
		Email:          credentials.Email,
		HashedPassword: string(hashedPassword),
	}

	// Next, insert into the dao
	err = h.dao.SaveUser(user)
	if err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
		return e.Wrap("error saving user to db")
	}

	return nil
}
