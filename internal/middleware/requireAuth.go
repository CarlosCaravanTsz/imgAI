package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CarlosCaravanTsz/imgAI/internal/auth"
	"github.com/CarlosCaravanTsz/imgAI/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims mínimos que esperas en el token
type MyClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func RequireAuth(c *gin.Context) {
	fmt.Println("In middleware")

	// Obtener la cookie de la request
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Obtener email desde los claims de la Cookie JWT 

	claims, err := auth.ExtractClaimsWithoutVerify(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Error: Invalid token"})
		return
	}

	email := claims["sub"].(string)

	// crear db conn
	db := database.GetConnection()

	// Obtener secret token del usuario
	var usuario database.Usuario
	if err := db.Where("email = ?", email).First(&usuario).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Error: Invalid credentials", "logged": false})
		return
	}

	// Decodificar/validar
	exp, err := auth.ValidateToken(tokenString, usuario.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Error: Invalid credentials", "logged": false})
		return
	}

	// Revisar la expiracion
	if float64(time.Now().Unix()) > exp {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "Error: Invalid credentials", "logged": false})
		return
	}

	// Agregar el userID a la req
	c.Set("user", usuario)

	c.Next()
}
