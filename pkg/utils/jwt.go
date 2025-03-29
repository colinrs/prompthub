package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(claimsInfo map[string]interface{}, signingKey []byte, expAt int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for k, v := range claimsInfo {
		claims[k] = v
	}
	claims["exp"] = float64(expAt)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string, signingKey []byte) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return true
		}
		return true
	}
	return false
}
