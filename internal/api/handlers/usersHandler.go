package handlers

import (
	"net/http"

	auth "github.com/CarlosCaravanTsz/imgAI/internal/auth"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	log "github.com/CarlosCaravanTsz/imgAI/internal/logger"

	"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type UsersRouteHandlers struct{}

func (h *UsersRouteHandlers) RegisterUser(c *gin.Context) {


	// Obtienes body values
	var newUser struct {
	Name   string `form:"name" json:"name"  binding:"required,max=100"`
	Email    string `form:"email" json:"email" binding:"required,email,max=100"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=64"`
}
	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hashear password
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(500, gin.H{"status": "error: Error hashing password"})
		return
	}

	// crear secret key para firma jwt: hashedPassword + randomGeneration + env secret
	s, err := auth.GenerateRandomString()
	if err != nil {		
		log.LogError("Error generating random string", logrus.Fields{
			"error": err,
		})
		c.JSON(500, gin.H{"status": "error: Error generating random string"})
		return
	}
	token, err := auth.GenerateToken(hashedPassword, s)
	if err != nil {
		log.LogError("Error generating user secret key", logrus.Fields{
			"error": err,
		})
		c.JSON(500, gin.H{"status":"error: Error generating user secret key"})
		return
	}

	db := database.GetConnection()

	usuario := database.Usuario{
		Nombre:       newUser.Name,
		Email:        newUser.Email,
		PasswordHash: hashedPassword,
		Token:        token,
	}

	// Creas usuario
	results := db.Create(&usuario)
	if results.Error != nil {
			c.JSON(500, gin.H{"error": results.Error})
		return
	}

	if results.Error != nil {
		c.JSON(500, gin.H{"error": results.Error})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "Usuario creado",
	})

}

func (h *UsersRouteHandlers) LoginUser(c *gin.Context) {

	var loginUser struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

	if err := c.ShouldBind(&loginUser); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := database.GetConnection()

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

	// Genera token
	token, err := auth.GenerateTokenJWT(usuario.Email, usuario.Token )
	if err != nil {
		c.JSON(401, gin.H{"status": "Error: Invalid to create token", "logged": false})
		return
	}

	// Setear cookie HTTP only
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600 * 24, "", "", true, true)

	c.JSON(201, gin.H{"status": "Logged correctly", "logged": true})

}

func GetUser(c *gin.Context) (*database.Usuario, bool) {
    u, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not set"})
        return nil, false
    }
    user, ok := u.(database.Usuario)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
        return nil, false
    }
    return &user, true
}

// func (h *UsersRouteHandlers) LogoutUser(c *gin.Context) {

 	// quitar el estado en el front
// }




