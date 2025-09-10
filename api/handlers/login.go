package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *APIHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signingKey := []byte("l9Xf7vQ3s8R2yT1uM6kB4hJ0aW5nZ8pC")
		body := make(map[string]string)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		login, login_ok := body["login"]
		password, password_ok := body["password"]
		db_password, db_ok := a.dbRepo.GetUser(login)
		if login_ok && password_ok && db_ok && db_password == password {
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
			claims["iat"] = time.Now().Unix()
			claims["sub"] = login

			tokenString, err := token.SignedString(signingKey)
			if err != nil {
				http.Error(w, "Error signing token", http.StatusInternalServerError)
				return
			}

			cookie := &http.Cookie{
				Name:     "jwt",
				Value:    tokenString,
				HttpOnly: true,
				Secure:   false,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
				SameSite: http.SameSiteLaxMode,
			}

			http.SetCookie(w, cookie)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Login successful"))
		} else {
			http.Error(w, "Wrong credentials", http.StatusBadRequest)
		}
	}
}
