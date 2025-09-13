package handlers

import (
	"Sugyk/jwt_golang/db_repository"
	"Sugyk/jwt_golang/jwt_repository"
	"net/http"
)

func Register(mux *http.ServeMux, dbRepo *db_repository.DBRepo, jwtRepo *jwt_repository.JWTRepo) {
	apiHandler := NewAPIHandler(dbRepo, jwtRepo)

	mux.Handle("/healthz", apiHandler.Health())
	mux.Handle("/reg", apiHandler.Register())
	mux.Handle("/login", apiHandler.Login())
	mux.Handle("/check_token", apiHandler.CheckJWT())
}
