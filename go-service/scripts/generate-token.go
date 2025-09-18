package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWT claims structure
type Claims struct {
	ID       string `json:"id"`
	CsrfHmac string `json:"csrf_hmac"`
	jwt.RegisteredClaims
}

func generateToken(id, csrfHmac, secret string, expiresIn time.Duration) (string, error) {
	claims := Claims{
		ID:       id,
		CsrfHmac: csrfHmac,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func generateCsrfHmacHash(csrfToken, secret string) string {
	key := []byte(secret)
	message := []byte(csrfToken)

	h := hmac.New(sha256.New, key)
	h.Write(message)

	return hex.EncodeToString(h.Sum(nil))
}

func Generate() {
	// Generate CSRF token
	csrfToken := uuid.New().String()
	fmt.Printf("CSRF Token: %s\n", csrfToken)

	// Generate CSRF HMAC hash
	csrfHmacHash := generateCsrfHmacHash(csrfToken, os.Getenv("JWT_ACCESS_TOKEN_SECRET"))
	fmt.Printf("CSRF HMAC Hash: %s\n", csrfHmacHash)

	// Generate access token
	accessToken, err := generateToken(
		"golang-service",
		csrfHmacHash,
		os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		time.Duration(15)*time.Minute, // Default 15 minutes, adjust as needed
	)
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}

	fmt.Printf("Access Token: %s\n", accessToken)
	fmt.Println("Use the above tokens in the Go service for authentication.")
}
