package handlers

import (
	"Sugyk/jwt_golang/blacklist_repository"
	"Sugyk/jwt_golang/db_repository"
	"net/http"
)

func Register(mux *http.ServeMux, dbRepo *db_repository.DBRepo, blRepo *blacklist_repository.BLRepo) {
	apiHandler := NewAPIHandler(dbRepo, blRepo)

	mux.Handle("/healthz", apiHandler.Health())
	mux.Handle("/reg", apiHandler.Register())
	mux.Handle("/login", apiHandler.Login())
	mux.Handle("/check_token", apiHandler.CheckJWT())
}
