package handlers

import (
	"net/http"
)

func parseJWTCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return "nil", err
	}
	return cookie.Value, nil
}
