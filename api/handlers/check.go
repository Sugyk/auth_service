package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrJWTMissing        = errors.New("JWT token is missing")
	ErrInvalidJWT        = errors.New("JWT token is invalid")
	ErrInvalidAuthHeader = errors.New("auth header is invalid")
	ErrDecodingJSON      = errors.New("error when decoding json")
)

func (a *APIHandler) CheckJWT() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth_header := strings.Split(r.Header.Get("Authorization"), " ")
		if len(auth_header) != 2 || auth_header[0] != "Bearer" {
			log.Printf("Authorization header has invalid part(s): %v", auth_header)
			http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
			return
		}
		jwtToken := auth_header[1]

		t, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
			return []byte(getSecretKey()), nil
		})

		if err != nil || !t.Valid {
			log.Println()
			log.Printf("error when validating jwt: %v", err)
			http.Error(w, ErrInvalidJWT.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("OK"))
	}
}
