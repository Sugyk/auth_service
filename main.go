package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var users_database = make(map[string]string)

func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]string)
		json.NewDecoder(r.Body).Decode(&body)

		login, login_ok := body["login"]
		password, password_ok := body["password"]
		_, db_ok := users_database[login]
		if login_ok && password_ok && !db_ok {
			users_database[login] = password
			w.Write([]byte("User created"))
		} else {
			w.Write([]byte("User is already appeares"))
		}
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signingKey := []byte("l9Xf7vQ3s8R2yT1uM6kB4hJ0aW5nZ8pC")
		body := make(map[string]string)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		login, login_ok := body["login"]
		password, password_ok := body["password"]
		db_password, db_ok := users_database[login]
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

func CheckJWT() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/healthz", Health())
	mux.Handle("/reg", Register())
	mux.Handle("/login", Login())
	mux.Handle("/check", CheckJWT())
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Error while listening:", err)
	}
}
