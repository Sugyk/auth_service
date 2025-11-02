package handlers

import (
	"Sugyk/jwt_golang/blacklist_repository"
	"Sugyk/jwt_golang/db_repository"
	"os"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	jwtSecret = []byte(secret)
}

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
