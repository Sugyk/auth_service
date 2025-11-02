package handlers

import (
	"Sugyk/jwt_golang/db_repository"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

var (
	ErrCreatingJWT = errors.New("error creating jwt")
)

func (a *APIHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]string)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		login, login_ok := body["login"]
		if !login_ok {
			http.Error(w, "Wrong credentials", http.StatusBadRequest)
			return
		}

		loginIsBlocked, err := a.blRepo.CheckAndIncrement(login)
		if err != nil {
			log.Println("error checking blocklist:", err)
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		if loginIsBlocked {
			http.Error(w, "too many tries, try again later", http.StatusTooManyRequests)
			return
		}
		password, password_ok := body["password"]
		hash_password, err := a.dbRepo.GetUserPasswordHash(login)
		if err != nil {
			if err == db_repository.ErrUserNotFound {
				http.Error(w, "Wrong credentials", http.StatusUnauthorized)
				return
			}
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		if password_ok {
			if !CheckPasswordHash(password, hash_password) {
				http.Error(w, "Wrong credentials", http.StatusUnauthorized)
				return
			}
			tokenString, err := CreateJWT(login)
			if err != nil {
				log.Printf("error creating jwt: %v", err)
				http.Error(w, ErrCreatingJWT.Error(), http.StatusInternalServerError)
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
