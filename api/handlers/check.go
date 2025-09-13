package handlers

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (a *APIHandler) CheckJWT() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := parseJWTCookie(r)
		if err != nil {
			log.Println("error when parsing token ", err)
		}

		t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			return []byte("l9Xf7vQ3s8R2yT1uM6kB4hJ0aW5nZ8pC"), nil
		})

		if err == nil && t.Valid {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		} else {
			http.Error(w, "Token is invalid or expired", http.StatusUnauthorized)
			return
		}
	}
}
