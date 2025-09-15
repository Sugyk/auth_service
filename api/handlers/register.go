package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func (a *APIHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]string)
		json.NewDecoder(r.Body).Decode(&body)

		login, login_ok := body["login"]
		password, password_ok := body["password"]
		hashedPassword, err := HashPassword(password)
		if err != nil {
			log.Println("error when hashing password")
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		if login_ok && password_ok {
			if err := a.dbRepo.CreateUser(login, hashedPassword); err == nil {
				w.Write([]byte("User created"))
				return
			}
		}
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("User with this login is already exists"))
	}
}
