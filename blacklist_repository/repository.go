package blacklist_repository

import "github.com/redis/go-redis/v9"

type BLRepo struct {
	db *redis.Client
}

func NewBLRepo(db *redis.Client) *BLRepo {
	return &BLRepo{
		db: db,
	}
}
