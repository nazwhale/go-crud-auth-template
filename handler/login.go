package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FilmListClub/backend/auth"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) *Error {
	var credentials Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		return &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusBadRequest,
		}
	}

	// Get the hashed password we're storing
	user, err := h.dao.ReadUserByEmail(credentials.Email)
	if err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
		return e.Wrap("error reading user from db")
	}

	// Compare the hashed password in the db with the password in the request body
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(credentials.Password)); err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusUnauthorized,
		}
		return e.Wrap("wrong password")
	}

	// Password is correct 🔑
	// Get a JWT token encrusted cookie
	cookie, err := auth.GetNewCookie(credentials.Email)
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
