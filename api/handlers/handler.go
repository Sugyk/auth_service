package handlers

import (
	"Sugyk/jwt_golang/db_repository"
)

type APIHandler struct {
	dbRepo *db_repository.DBRepo
}

func NewAPIHandler(dbRepo *db_repository.DBRepo) *APIHandler {
	return &APIHandler{
		dbRepo: dbRepo,
	}
}
