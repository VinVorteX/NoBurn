package utils

import (
	"testing"
	"time"

	"github.com/VinVorteX/NoBurn/internal/config"
)

func TestGenerateToken(t *testing.T) {
	config.AppConfig = &config.Config{
		JwtSecret:    "test-secret-key",
		JwtExpiresIn: 24 * time.Hour,
	}

	token, err := GenerateToken(1, "test@example.com")
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("Expected non-empty token")
	}
}

func TestVerifyToken(t *testing.T) {
	config.AppConfig = &config.Config{
		JwtSecret:    "test-secret-key",
		JwtExpiresIn: 24 * time.Hour,
	}

	token, _ := GenerateToken(1, "test@example.com")

	parsedToken, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	if !parsedToken.Valid {
		t.Error("Expected valid token")
	}
}

func TestVerifyInvalidToken(t *testing.T) {
	config.AppConfig = &config.Config{
		JwtSecret: "test-secret-key",
	}

	_, err := VerifyToken("invalid.token.here")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}