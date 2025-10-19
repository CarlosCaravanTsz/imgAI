package main

import (

	"github.com/CarlosCaravanTsz/imgAI/internal/api"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	l "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		l.LogInfo("Error while uploading .env vars", logrus.Fields{
			"error": err,
		})
	}
	database.Connect()
}

func main() {
	r := gin.Default()
	api.RegisterRoutes(r)

	// r.LoadHTMLGlob("templates/*")
	// r.Static("/static", "./templates") // Add this line to serve static files
	

	r.Run(":8080")
}