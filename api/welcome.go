package api

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	// Obtain the session token from the requests cookies
	c, err := r.Cookie(cookieNameToken)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := c.Value
	claims := &Claims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}

	// Parse the JWT string and store the result in `claims`.
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("👋 Welcome: %s", claims.Username)))
}
