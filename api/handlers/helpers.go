package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func getSecretKey() []byte {
	return jwtSecret
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DecodeBody(w http.ResponseWriter, r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return err
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("password must be at least 12 characters long")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func CreateJWT(login string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["iat"] = time.Now().Unix()
	claims["sub"] = login

	tokenString, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
