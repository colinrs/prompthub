package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HashPassword(password, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(password))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	return hashedPassword
}

func CheckPassword(password, hashedPassword, key string) bool {
	hashedInput := HashPassword(password, key)
	return hmac.Equal([]byte(hashedInput), []byte(hashedPassword))
}
