package jwt_repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (j *JWTRepo) GetJWT(login string) (string, error) {
	jwtstring, err := j.db.Get(context.Background(), login).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return jwtstring, err
}
