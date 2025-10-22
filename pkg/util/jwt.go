package util

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT() {}

func GenerateJWT(userId, secret, issuer string, expiration time.Duration) (string, error) {
	secretKey := []byte(secret)
	expirationTime := time.Now().Add(expiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expirationTime,
		"iat":    time.Now().Unix(),
		"iss":    issuer,
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("error signing token: %v", err)
		return "", err
	}
	return tokenString, nil
}
