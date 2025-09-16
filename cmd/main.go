package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"Sugyk/jwt_golang/blacklist_repository"
	"Sugyk/jwt_golang/db_repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/redis/go-redis/v9"

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
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not up migrates: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("No changes to migrate")
	} else {
		log.Println("Migrations applied successfully")
	}

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

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("redis connection error: %w", err)
	}
	blRepo := blacklist_repository.NewBLRepo(redisClient)

	handlers.Register(mux, dbRepo, blRepo)
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Println("Error while listening:", err)
	}
}
