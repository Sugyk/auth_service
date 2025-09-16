package handlers

import (
	"Sugyk/jwt_golang/blacklist_repository"
	"Sugyk/jwt_golang/db_repository"
)

type APIHandler struct {
	dbRepo *db_repository.DBRepo
	blRepo *blacklist_repository.BLRepo
}

func NewAPIHandler(dbRepo *db_repository.DBRepo, blRepo *blacklist_repository.BLRepo) *APIHandler {
	return &APIHandler{
		dbRepo: dbRepo,
		blRepo: blRepo,
	}
}
