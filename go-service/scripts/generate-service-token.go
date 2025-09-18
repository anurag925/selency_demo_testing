package main

import (
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

func generateServiceToken() (string, error) {
	// Implementation for generating a service token
	jwtSecret := os.Getenv("SERVICE_TOKEN_SECRET")
	claims := jwt.MapClaims{
		"service": "golang-service",
	}
	slog.Info("secret", jwtSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringtoken, err := token.SignedString([]byte(jwtSecret))
	slog.Info("Generated service token", "token", stringtoken)
	return stringtoken, err
}
