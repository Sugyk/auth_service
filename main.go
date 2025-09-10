package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"log"
	"net/http"
)

var users_database = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	handlers.Register(mux, users_database)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Error while listening:", err)
	}
}
