package handlers

import "net/http"

func Register(mux *http.ServeMux, users_database map[string]string) {
	apiHandler := APIHandler{
		users_database: users_database,
	}
	mux.Handle("/healthz", apiHandler.Health())
	mux.Handle("/reg", apiHandler.Register())
	mux.Handle("/login", apiHandler.Login())
}
