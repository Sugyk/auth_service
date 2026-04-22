package jwt_manager

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	CreateJWT(login string) (string, error)
}

type jwtManager struct {
	secret []byte
	ttl    time.Duration
}

func (j *jwtManager) CreateJWT(login string) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": login,
		"iat": now.Unix(),
		"exp": now.Add(j.ttl).Unix(),
	})

	tokenString, err := token.SignedString(j.secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewJWTManager(secret []byte, ttl time.Duration) (JWTManager, error) {
	if secret == nil {
		return nil, fmt.Errorf("secret for JWTManager can not be empty")
	}

	if ttl <= 0 {
		return nil, fmt.Errorf("ttl for token must be greater than 0")
	}

	return &jwtManager{
		secret: secret,
		ttl:    ttl,
	}, nil
}
