package handlers

import (
	"Sugyk/jwt_golang/db_repository"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (a *APIHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signingKey := []byte(getSecretKey())
		body := make(map[string]string)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		login, login_ok := body["login"]
		password, password_ok := body["password"]
		hash_password, err := a.dbRepo.GetUserPasswordHash(login)
		if err != nil {
			if err == db_repository.ErrUserNotFound {
				http.Error(w, "Wrong credentials", http.StatusBadRequest)
				return
			}
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		if login_ok && password_ok {
			if !CheckPasswordHash(password, hash_password) {
				http.Error(w, "Wrong credentials", http.StatusBadRequest)
			}
			token := jwt.New(jwt.SigningMethodHS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
			claims["iat"] = time.Now().Unix()
			claims["sub"] = login

			tokenString, err := token.SignedString(signingKey)
			if err != nil {
				log.Println("Error signing token")
				http.Error(w, "something went wrong", http.StatusInternalServerError)
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
