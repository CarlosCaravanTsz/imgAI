package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenJWT(userID uint, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (float64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return  0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["exp"].(float64), nil
	} else {
		return  0, err
	}
}
