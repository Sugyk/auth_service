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
		_, db_ok := a.users_database[login]
		if login_ok && password_ok && !db_ok {
			a.users_database[login] = password
			w.Write([]byte("User created"))
		} else {
			w.Write([]byte("User is already appeares"))
		}
	}
}
