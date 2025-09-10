package handlers

import (
	"encoding/json"
	"net/http"
)

func (a *APIHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := make(map[string]string)
		json.NewDecoder(r.Body).Decode(&body)

		login, login_ok := body["login"]
		password, password_ok := body["password"]

		if login_ok && password_ok {
			if err := a.dbRepo.CreateUser(login, password); err == nil {
				w.Write([]byte("User created"))
				return
			}
		}
		w.Write([]byte("User is already appeares"))
	}
}
