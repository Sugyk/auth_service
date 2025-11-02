package main

import (
	"Sugyk/jwt_golang/api/handlers"
	"Sugyk/jwt_golang/blacklist_repository"
	"Sugyk/jwt_golang/db_repository"
	"Sugyk/jwt_golang/packages/configs"
	"Sugyk/jwt_golang/packages/database"
	"Sugyk/jwt_golang/packages/migrations"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	mux := http.NewServeMux()
	db_config := configs.NewDBConfig()
	db, err := database.NewDbConnection(
		db_config.Username,
		db_config.Password,
		db_config.Address,
		db_config.Port,
		db_config.Db_name,
		db_config.Ssl_mode,
	)
	if err != nil {
		log.Fatalf("error when opening db connection: %v", err)
	}
	defer db.Close()
	log.Println("db connection is established")
	err = migrations.RunMigrations(db)
	if err != nil {
		log.Fatal("migration failed:", err)
	}

	dbRepo := db_repository.NewDBRepo(db)
	redisConfig, err := configs.NewRedisConfig()
	if err != nil {
		log.Fatalf("error getting configs for redis: %v", err)
	}

	redisClient, err := database.NewRedisConnection(
		redisConfig.Addr,
		redisConfig.Password,
		redisConfig.DB,
	)
	if err != nil {
		log.Fatalf("error connecting redis: %v", err)
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
