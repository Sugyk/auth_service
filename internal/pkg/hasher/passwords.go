package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type passwordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) *passwordHasher {
	return &passwordHasher{
		cost: cost,
	}
}

func (p *passwordHasher) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)

	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (p *passwordHasher) CompareHashAndPassword(password string, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}
	return true
}
