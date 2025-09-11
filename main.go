package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"Sugyk/jwt_golang/db_repository"
	"log"
	"net/http"
)

var users_database = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	dbRepo := db_repository.NewDBRepo()
	handlers.Register(mux, dbRepo)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	redis_client, err := NewClient(Config{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	if err != nil {
		panic(err)
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error while listening:", err)
	}
}
