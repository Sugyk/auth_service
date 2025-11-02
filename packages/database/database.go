package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisConnection(
	address string,
	password string,
	db_number int,
) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db_number,
	})
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return redisClient, nil
}

func NewDbConnection(
	username string,
	password string,
	address string,
	port string,
	db_name string,
	ssl_mode string,
) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		username,
		password,
		address,
		port,
		db_name,
		ssl_mode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
