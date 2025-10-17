package main

import (
	"github.com/CarlosCaravanTsz/imgAI/internal/api"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	l "github.com/CarlosCaravanTsz/imgAI/internal/logger"


)

func init() {
	if err := godotenv.Load(); err != nil {
		l.LogInfo("Error while uploading ENV vars", logrus.Fields{
			"error": err,
		})
}
}

func main() {

	database.Connect()
	
	r := gin.Default()

	//r.LoadHTMLGlob("templates/*")
	//r.Static("/static", "./templates") // Add this line to serve static files

	api.RegisterRoutes(r)

	r.Run(":8081")
}
