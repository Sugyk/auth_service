package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"Sugyk/jwt_golang/db_repository"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	mux := http.NewServeMux()
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	dbRepo := db_repository.NewDBRepo()

	handlers.Register(mux, dbRepo)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error while listening:", err)
	}
}
