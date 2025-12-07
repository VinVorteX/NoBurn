package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/VinVorteX/NoBurn/internal/config"
)

func GenerateToken(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(config.AppConfig.JwtExpiresIn).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JwtSecret))
}

func VerifyToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JwtSecret), nil
	})
}