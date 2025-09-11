package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"Sugyk/jwt_golang/db_repository"
	"Sugyk/jwt_golang/jwt_repository"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func main() {
	mux := http.NewServeMux()

	dbRepo := db_repository.NewDBRepo()
	jwtClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	JWTRepo := jwt_repository.NewJWTRepo(jwtClient)

	handlers.Register(mux, dbRepo, JWTRepo)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error while listening:", err)
	}
}
