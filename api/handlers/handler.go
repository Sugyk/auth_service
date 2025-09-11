package handlers

import (
	"Sugyk/jwt_golang/db_repository"
	"Sugyk/jwt_golang/jwt_repository"
)

type APIHandler struct {
	dbRepo  *db_repository.DBRepo
	jwtRepo *jwt_repository.JWTRepo
}

func NewAPIHandler(dbRepo *db_repository.DBRepo, jwtRepo *jwt_repository.JWTRepo) *APIHandler {
	return &APIHandler{
		dbRepo:  dbRepo,
		jwtRepo: jwtRepo,
	}
}
