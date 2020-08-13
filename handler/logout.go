package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/FilmListClub/backend/auth"
)

type logoutReq struct {
	ID int `json:"id"`
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) *Error {
	// TODO: pull CORS out into middleware
	(w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")

	req := &logoutReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		e := &Error{
			Message:      err.Error(),
			ResponseCode: http.StatusBadRequest,
		}
		return e.Wrap("error decoding json body")
	}

	switch {
	case !auth.IsIDValid(req.ID):
		return &Error{
			Message:      "invalid ID",
			ResponseCode: http.StatusBadRequest,
		}
	}

	c, err := r.Cookie(auth.GetCookieNameToken())
	if err != nil {
		fmt.Println("cook:", c)
		switch {
		// If there's no cookie we're good, as we'd only delete it anyway
		case err == http.ErrNoCookie:
			fmt.Println("No cookie in the first place")
			return nil
		default:
			e := &Error{
				Message:      err.Error(),
				ResponseCode: http.StatusBadRequest,
			}
			return e.Wrap("cookie err")
		}
	}

	// Set maxAge to 0 to make it expired
	c.Expires = time.Now().UTC()

	//// To delete the old cookie, we set a new one with a maxAge of 0
	//cookie := auth.GetOldCookie()

	http.SetCookie(w, c)

	fmt.Println("Set cookie with expiry of now", c)

	return nil
}
