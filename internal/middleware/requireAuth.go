package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CarlosCaravanTsz/imgAI/internal/auth"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)



func RequireAuth(c *gin.Context) {
	fmt.Println("In middleware")

	// Obtener la cookie de la request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// crear db conn
	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	var loginUser struct {
			Email    string `form:"email" binding:"required,email,max=100"`
	}

	if err := c.ShouldBind(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

		// Obtener el id y el token del usuario
		var usuario database.Usuario
	if err := db.Where("email = ?", loginUser.Email).First(&usuario).Error; err != nil {
		c.JSON(401, gin.H{"status": "Error: Invalid credentials", "logged": false})
		return
	}


	// Decodificar/validar

	exp, err := auth.ValidateToken(tokenString, usuario.Token)

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// Revisar la expiracion

	if float64(time.Now().Unix()) > exp {
		c.AbortWithStatus(http.StatusUnauthorized)
	}


	// Agregar el userID a la req
	c.Set("user", usuario)

	c.Next()
}
