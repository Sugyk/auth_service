package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (a *APIHandler) CheckJWT() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token, err := parseJWTCookie(r)
		if err != nil {
			http.Error(w, "Token is invalid or expired", http.StatusBadRequest)
			return
		}

		t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
			return []byte(getSecretKey()), nil
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
