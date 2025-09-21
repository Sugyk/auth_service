package handlers

import (
	"Sugyk/jwt_golang/api/middlewares"
	"Sugyk/jwt_golang/blacklist_repository"
	"Sugyk/jwt_golang/db_repository"
	"net/http"
)

func Register(mux *http.ServeMux, dbRepo *db_repository.DBRepo, blRepo *blacklist_repository.BLRepo) {
	apiHandler := NewAPIHandler(dbRepo, blRepo)

	mux.Handle("/healthz", middlewares.Get(apiHandler.Health()))
	mux.Handle("/reg", middlewares.Post(apiHandler.Register()))
	mux.Handle("/login", middlewares.Post(apiHandler.Login()))
	mux.Handle("/check_token", middlewares.Get(apiHandler.CheckJWT()))
}
