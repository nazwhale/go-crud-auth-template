package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

const cookieNameToken = "token"
const cookieExpiration = 5 * time.Minute

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GetCookieNameToken() string {
	return cookieNameToken
}

func GetNewCookie(email string) (*http.Cookie, error) {
	expirationTime := time.Now().Add(cookieExpiration)

	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token, err := newJWTFromClaims(claims)
	if err != nil {
		return nil, errors.New("error creating JWT string")
	}

	return &http.Cookie{
		Name:    cookieNameToken,
		Value:   token,
		Expires: expirationTime,
	}, nil
}

func GetOldCookie() *http.Cookie {
	// To delete a cookie, we set the max age to 0
	return &http.Cookie{
		Name:   cookieNameToken,
		MaxAge: 0,
	}
}

func GetRefreshedCookie(claims *Claims) (*http.Cookie, error) {
	expirationTime := time.Now().Add(cookieExpiration)

	// Refresh expiration time
	claims.ExpiresAt = expirationTime.Unix()

	token, err := newJWTFromClaims(claims)
	if err != nil {
		return nil, errors.New("error creating JWT string")
	}

	return &http.Cookie{
		Name:    cookieNameToken,
		Value:   token,
		Expires: expirationTime,
	}, nil
}

func newJWTFromClaims(claims *Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateRequest(r *http.Request) (*Claims, error) {
	c, err := r.Cookie(cookieNameToken)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, errors.New("no cookie")
		}

		return nil, errors.New("could not validate request")
	}

	tokenString := c.Value
	claims := &Claims{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	}

	// Parse the JWT string and store the result in claims
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("jwt signature invalid")
		}

		return nil, errors.New("jwt parsing error")
	}
	if !token.Valid {
		return nil, errors.New("jwt invalid")
	}

	return claims, nil
}
