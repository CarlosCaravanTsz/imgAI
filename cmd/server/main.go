package main

<<<<<<< HEAD
import (
	"github.com/CarlosCaravanTsz/imgAI/internal/api"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()
	
	r := gin.Default()

	//r.LoadHTMLGlob("templates/*")
	//r.Static("/static", "./templates") // Add this line to serve static files

	api.RegisterRoutes(r)

	r.Run(":8081")
=======
func main() {
>>>>>>> frontend
}
