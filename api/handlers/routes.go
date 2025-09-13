package handlers

import (
	"Sugyk/jwt_golang/db_repository"
	"net/http"
)

func Register(mux *http.ServeMux, dbRepo *db_repository.DBRepo) {
	apiHandler := NewAPIHandler(dbRepo)

	mux.Handle("/healthz", apiHandler.Health())
	mux.Handle("/reg", apiHandler.Register())
	mux.Handle("/login", apiHandler.Login())
	mux.Handle("/check_token", apiHandler.CheckJWT())
}
