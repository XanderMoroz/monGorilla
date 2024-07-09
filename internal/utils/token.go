package utils

import (
	"crypto/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func NewID() string {
	b := make([]byte, 20)
	rand.Read(b)
	return string(b)
}

func CreateJSONWebToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        NewID(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	token_string, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return token_string, nil
}
