package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"Sugyk/jwt_golang/db_repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connection is established")
	dbRepo := db_repository.NewDBRepo(db)

	handlers.Register(mux, dbRepo)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error while listening:", err)
	}
}
