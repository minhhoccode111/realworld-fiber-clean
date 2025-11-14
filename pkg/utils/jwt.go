package utils

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID, secret, issuer string, expiration time.Duration) (string, error) {
	secretKey := []byte(secret)
	expirationTime := time.Now().Add(expiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID, // NOTE: use only 'userID' is insecure?
		"exp": expirationTime,
		"iat": time.Now().Unix(),
		"iss": issuer,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Printf("error signing token: %v", err)

		return "", err
	}

	return tokenString, nil
}
