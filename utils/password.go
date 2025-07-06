package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given plain-text password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a hashed password with a plain-text one.
func CheckPasswordHash(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateSalt generates a random 16-byte salt and returns it as a base64 string.
func GenerateSalt() (string, error) {
	saltBytes := make([]byte, 16)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return base64.StdEncoding.EncodeToString(saltBytes), nil
}

// HashPasswordWithSalt hashes the password combined with a salt.
func HashPasswordWithSalt(password, salt string) (string, error) {
	combined := password + salt
	bytes, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHashWithSalt checks the password using stored hash and salt.
func CheckPasswordHashWithSalt(hashedPassword, password, salt string) bool {
	combined := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(combined))
	return err == nil
}
