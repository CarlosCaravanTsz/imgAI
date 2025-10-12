package auth

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/joho/godotenv"
)

func GenerateRandomString() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func GenerateToken(hashedPassword, randomString string) (string, error) {
	err := godotenv.Load() // dsn := os.Getenv("DATABASE_URL")
	if err != nil {
		log.LogError("Error loading .env file in auth", logrus.Fields{
			"error": err,
		})
		return "", err
	}

	jwt := os.Getenv("SECRET_JWT")

	return hashedPassword + randomString + jwt, nil
}
