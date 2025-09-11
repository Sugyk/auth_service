package jwt_repository

import (
	"context"
	"time"
)

func (j *JWTRepo) SetJWT(login string, jwt string) error {
	status := j.db.Set(context.Background(), login, jwt, time.Minute*15)
	if err := status.Err(); err != nil {
		return err
	}
	return nil
}
