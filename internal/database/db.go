package database

import (
    "os"
    "sync"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    l "github.com/CarlosCaravanTsz/imgAI/internal/logger"
    "github.com/sirupsen/logrus"
)

var (
    DB   *gorm.DB
    once sync.Once
)

func Connect() *gorm.DB {
    once.Do(func() {
        dsn := os.Getenv("DATABASE_URL")
        db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
        if err != nil {
            l.LogError("Error connecting to DB", logrus.Fields{"error": err})
            panic(err)
        }
        DB = db
        if err := DB.AutoMigrate(&Usuario{}, &Foto{}, &Album{}); err != nil {
            l.LogError("Error during migrations", logrus.Fields{"error": err})
        }
    })
    return DB
}

func GetConnection() *gorm.DB {
    if DB == nil {
        return Connect()
    }
    return DB
}
