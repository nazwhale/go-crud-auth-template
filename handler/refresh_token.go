package handler

import (
	"net/http"
	"time"

	"github.com/FilmListClub/backend/auth"
)

const expiryThreshold = 30 * time.Second

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) *Error {
	claims, err := auth.ValidateRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusBadRequest,
		}
		return e.Wrap("error validating request")
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > expiryThreshold {
		return &Error{
			Message:      "token expiry not within threshold",
			ResponseCode: http.StatusBadRequest,
		}
	}

	// Create a new token with a renewed expiration time
	cookie, err := auth.GetRefreshedCookie(claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusInternalServerError,
		}
		return e.Wrap("error refreshing cookie")
	}

	http.SetCookie(w, cookie)

	return nil
}
