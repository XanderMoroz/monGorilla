package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewID() string {
	b := make([]byte, 20)
	rand.Read(b)
	return string(b)
}

func CreateJSONWebToken(userId primitive.ObjectID) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": userId,
		"issued":  jwt.NewNumericDate(time.Now()),
		"expires": jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Expires in 24 hours
	})

	token_string, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return token_string, nil
}

func ParseUserIDFromJWTToken(cookieWithJWT string) (string, error) {

	api_secret := os.Getenv("JWT_SECRET_KEY")
	hmacSecret := []byte(api_secret)

	token, err := jwt.Parse(cookieWithJWT, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена и другие параметры
		return hmacSecret, nil
	})

	if err != nil {
		// log.Printf("При извлечении токена произошла ошибка <%v>\n", err)
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Неверный JWT токен")
		return "", fmt.Errorf("invalid JWT Token ")
	}

	userID, ok := claims["subject"].(string)
	if !ok {
		log.Println("Не удалось извлечь USER_ID из токена")
		return "", fmt.Errorf("failed to parse claims")
	}

	return userID, nil
}
