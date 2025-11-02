package configs

import (
	"os"
	"strconv"
)

type DbConnectionConfig struct {
	Username string
	Password string
	Address  string
	Port     string
	Db_name  string
	Ssl_mode string
}

func NewDBConfig() *DbConnectionConfig {
	return &DbConnectionConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Address:  os.Getenv("DB_ADDRESS"),
		Port:     os.Getenv("DB_PORT"),
		Db_name:  os.Getenv("DB_NAME"),
		Ssl_mode: os.Getenv("DB_SSLMODE"),
	}
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisConfig() (*RedisConfig, error) {
	db_number, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}
	return &RedisConfig{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db_number,
	}, nil
}
