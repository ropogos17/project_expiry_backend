package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argon2Memory      uint32 = 19 * 1024 // 19 MiB
	argon2Iterations  uint32 = 2
	argon2Parallelism uint8  = 1
	argon2SaltLength  uint32 = 16
	argon2KeyLength   uint32 = 32
)

var (
	ErrInvalidPasswordHash = errors.New("invalid password hash")
	ErrPasswordMismatch    = errors.New("password mismatch")
)

// HashPassword generates an Argon2id password hash.
func HashPassword(password string) (string, error) {
	salt := make([]byte, argon2SaltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Iterations,
		argon2Memory,
		argon2Parallelism,
		argon2KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// PHC string format
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		argon2Memory,
		argon2Iterations,
		argon2Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

// CheckPassword checks if the provided password is correct or not
func CheckPassword(password string, encodedHash string) error {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" {
		return ErrInvalidPasswordHash
	}

	version, err := strconv.ParseUint(
		strings.TrimPrefix(parts[2], "v="), 10, 32,
	)
	if err != nil || uint32(version) != argon2.Version ||
		!strings.HasPrefix(parts[2], "v=") {
		return ErrInvalidPasswordHash
	}

	var memory, iterations uint32
	var parallelism uint8

	seen := make(map[string]bool)

	for _, param := range strings.Split(parts[3], ",") {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) != 2 || seen[kv[0]] {
			return ErrInvalidPasswordHash
		}
		seen[kv[0]] = true

		switch kv[0] {
		case "m":
			value, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return ErrInvalidPasswordHash
			}
			memory = uint32(value)

		case "t":
			value, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return ErrInvalidPasswordHash
			}
			iterations = uint32(value)

		case "p":
			value, err := strconv.ParseUint(kv[1], 10, 8)
			if err != nil {
				return ErrInvalidPasswordHash
			}
			parallelism = uint8(value)

		default:
			return ErrInvalidPasswordHash
		}
	}

	if len(seen) != 3 ||
		memory < argon2Memory || memory > 64*1024 ||
		iterations < argon2Iterations || iterations > 10 ||
		parallelism == 0 || parallelism > 8 ||
		memory < 8*uint32(parallelism) {
		return ErrInvalidPasswordHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < 16 || len(salt) > 64 {
		return ErrInvalidPasswordHash
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expectedHash) < 16 || len(expectedHash) > 64 {
		return ErrInvalidPasswordHash
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(expectedHash)),
	)

	if subtle.ConstantTimeCompare(actualHash, expectedHash) != 1 {
		return ErrPasswordMismatch
	}

	return nil
}
