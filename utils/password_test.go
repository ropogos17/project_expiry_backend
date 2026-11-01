package utils

import (
	"strings"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "MyStrongPassword123!"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Fatal("password hash is not Argon2id")
	}

	if err := CheckPassword(password, hash); err != nil {
		t.Fatalf("password verification failed: %v", err)
	}
}

func TestCheckPasswordWrongPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}

	err = CheckPassword("wrong-password", hash)

	if err == nil {
		t.Fatal("expected password verification to fail")
	}
}

func TestHashPasswordUniqueSalt(t *testing.T) {
	password := "same-password"

	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	if hash1 == hash2 {
		t.Fatal("expected different hashes due to random salts")
	}
}
