package utils

import "testing"

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"
	
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	if hash == password {
		t.Error("Hash should not equal plain password")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mySecurePassword123"
	hash, _ := HashPassword(password)

	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}

	if CheckPassword("wrongPassword", hash) {
		t.Error("CheckPassword should return false for incorrect password")
	}
}