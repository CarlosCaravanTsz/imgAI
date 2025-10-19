package auth

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"crypto/sha256"
	"crypto/hmac"

)

func GenerateRandomString() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateToken(hashedPassword, randomString string) (string, error) {
    jwt := os.Getenv("SECRET_JWT")
    key := []byte(jwt)

    mac := hmac.New(sha256.New, key)
    mac.Write([]byte(hashedPassword + randomString))
    return hex.EncodeToString(mac.Sum(nil)), nil
}