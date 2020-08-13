package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FilmListClub/backend/dao"

	"github.com/FilmListClub/backend/auth"

	"golang.org/x/crypto/bcrypt"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) *Error {
	// TODO: pull CORS out into middleware
	(w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")

	req := &loginReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusBadRequest,
		}
		return e.Wrap("error decoding json body")
	}

	switch {
	case !auth.IsEmailValid(req.Email):
		return &Error{
			Message:      "invalid email address",
			ResponseCode: http.StatusBadRequest,
		}
	case len(req.Password) < auth.MinimumPasswordLength:
		return &Error{
			Message:      fmt.Sprintf("password must be at least %d characters", auth.MinimumPasswordLength),
			ResponseCode: http.StatusBadRequest,
		}
	}

	// Get the hashed password we're storing
	user, err := h.dao.ReadUserByEmail(req.Email)
	if err != nil {
		e := &Error{
			Message: err.Error(),
		}

		switch err {
		case dao.ErrNoUserExists:
			e.ResponseCode = http.StatusUnauthorized
		default:
			e.ResponseCode = http.StatusInternalServerError
		}

		return e.Wrap("error reading user from db")
	}

	// Compare the hashed password in the db with the password in the request body
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusUnauthorized,
		}
		return e.Wrap("wrong password")
	}

	// Password is correct 🔑
	// Get a JWT token encrusted cookie
	cookie, err := auth.GetNewCookie(req.Email)
	if err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
		return e.Wrap("error creating cookie")
	}

	http.SetCookie(w, cookie)

	return nil
}
