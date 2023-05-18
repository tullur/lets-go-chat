// Package hasher provides methods to generate hash for the given passwords
// using a secure hashing algorithm - argon2
package hasher

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2 configuration parameters
const (
	iterations  uint32 = 3
	memory      uint32 = 64 * 1024
	parallelism uint8  = 2
	keyLength   uint32 = 32
)

// HashPassword generates a secure Argon2 password hash
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := generateArgonIdKey(password, salt)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedPassword := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, memory, iterations, parallelism, b64Salt, b64Hash,
	)

	return encodedPassword, nil
}

// CheckPasswordHash is used to compare hashed password with a plain-text password
func CheckPasswordHash(password, hash string) bool {
	hashParts := strings.Split(hash, "$")

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return false
	}

	hashToCompare := generateArgonIdKey(password, salt)

	return subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1
}

func generateArgonIdKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)
}
