package Middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
)

func SplitToken(headerToken string) string {
	parsToken := strings.SplitAfter(headerToken, " ")
	return parsToken[1]
}

func AuthenticateToken(tokenString string) error {
	_, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return err
	}

	return nil
}
