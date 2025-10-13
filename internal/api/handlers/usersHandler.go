package handlers

import (
	"fmt"
	"net/http"

	auth "github.com/CarlosCaravanTsz/imgAI/internal/auth"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"

	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type NewUserForm struct {
	Nombre   string `form:"nombre" binding:"required,max=100"`
	Email    string `form:"email" binding:"required,email,max=100"`
	Password string `form:"password" binding:"required,min=8,max=64"`
}

type LoginUser struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type UsersRouteHandlers struct{}

func (h *UsersRouteHandlers) RegisterUser(c *gin.Context) {
	nombre := c.PostForm("nombre")
	email := c.PostForm("email")
	password := c.PostForm("password")

	var newUser NewUserForm

	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Print(newUser)

	// hashear password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error hashing password"})
	}

	// crear token hashedPassword + randomGeneration + env secret
	s, _ := auth.GenerateRandomString()
	token, err := auth.GenerateToken(hashedPassword, s)
	if err != nil {
		log.LogError("Error loading .env file in auth", logrus.Fields{
			"error": err,
		})
	}

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	usuario := database.Usuario{
		Nombre:       nombre,
		Email:        email,
		PasswordHash: hashedPassword,
		Token:        token,
	}

	results := db.Create(&usuario)

	if results.Error != nil {
		c.JSON(500, gin.H{"error": results.Error})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "Usuario creado",
	})

}

func (h *UsersRouteHandlers) LoginUser(c *gin.Context) {
	var loginUser LoginUser

	if err := c.ShouldBind(&loginUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err := database.GetConnection()
	if err != nil {
		log.LogInfo("Error getting the db conn", logrus.Fields{})
	}

	// Find the user by email
	var usuario database.Usuario
	if err := db.Where("email = ?", loginUser.Email).First(&usuario).Error; err != nil {
		c.JSON(401, gin.H{"status": "Error: Invalid credentials", "logged": false})
		return
	}

	// Check Password
	if !auth.CheckPasswordHash(loginUser.Password, usuario.PasswordHash) {
		c.JSON(401, gin.H{"status": "Error: Invalid credentials", "logged": false})
		return
	}

	// TODO If password correct -> get the token and create a JWT
	c.JSON(201, gin.H{"status": "Logged correctly", "logged": true})

}

func (h *UsersRouteHandlers) LogoutUser(c *gin.Context) {

	// quitar el estado en el front
}
