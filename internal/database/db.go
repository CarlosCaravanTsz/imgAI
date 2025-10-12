package database

import (
	"log"
	"os"
	_ "os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	l "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConnection() (*gorm.DB, error) {
	err := godotenv.Load() // dsn := os.Getenv("DATABASE_URL")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		l.LogError("Error loading .env file", logrus.Fields{
			"error": err,
		})
		return nil, err
	}

	return db, nil
}

func Connect() {
	db, _ := GetConnection()

	err := db.AutoMigrate(Usuario{}, Foto{}, Album{})
	if err != nil {
		l.LogError("Error while doing migrations", logrus.Fields{
			"error": err,
		})
	}
}
