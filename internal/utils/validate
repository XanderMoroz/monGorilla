package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
)

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		flagUpper, flagLower, flagDigit, flagSpecial, flagSpace bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			flagUpper = true
		case unicode.IsLower(char):
			flagLower = true
		case unicode.IsSpace(char):
			flagSpace = true
		case unicode.IsDigit(char):
			flagDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			flagSpecial = true
		}
	}

	return flagUpper && flagLower && flagDigit && flagSpecial && !flagSpace
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	secret_key := os.Getenv("JWT_SECRET_KEY")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret_key), nil
	})
}

func IsStringEqual(args1, args2 string) bool {
	return args1 == args2
}

func ValidateCheckSpaceCharacter(args ...string) bool {
	flag := true

	l := len(args)
	wg := sync.WaitGroup{}
	wg.Add(l)

	lock := sync.Mutex{}

	for _, val := range args {
		go func(val string) {
			defer wg.Done()
			lock.Lock()
			if strings.TrimSpace(val) != val || val == "" {
				flag = false
			}
			lock.Unlock()
		}(val)
	}

	wg.Wait()

	return flag
}
