package jwt_repository

import "github.com/redis/go-redis/v9"

type JWTRepo struct {
	db *redis.Client
}

func NewJWTRepo(redisClient *redis.Client) *JWTRepo {
	return &JWTRepo{
		db: redisClient,
	}
}
