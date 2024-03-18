package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Error occurred while hashing the password: %v", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), []byte(password),
	)
	if err != nil {
		t.Error("Hashed password does not match the expected hash")
	}
}

func TestComparePasswordAndHash(t *testing.T) {
	password := "password"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	result := ComparePasswordAndHash(password, string(hash))
	if !result {
		t.Error("Hashed password does not compare")
	}
}
