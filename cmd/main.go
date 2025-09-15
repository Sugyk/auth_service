package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"Sugyk/jwt_golang/db_repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("could not up migrates: %w", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}

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
		log.Fatal("error when opening db connection:", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("db connection is established")
	err = runMigrations(db)
	if err != nil {
		log.Fatal("migration failed:", err)
	}

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
