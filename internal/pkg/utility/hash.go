package utility

import (
	"crypto/rand"
	"crypto/subtle"

	"golang.org/x/crypto/argon2"
)

const (
	iteration    = 1
	memory       = 64 * 2024
	parrallelism = 1
	keyLength    = 32
)

func HashPassword(password string) ([]byte, []byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	hash := argon2.IDKey([]byte(password), salt, iteration, memory, parrallelism, keyLength)
	return hash, salt, nil
}

func VerifyPassword(password string, salt, storedHash []byte) bool {
	hash := argon2.IDKey([]byte(password), salt, iteration, memory, parrallelism, keyLength)
	return subtle.ConstantTimeCompare(storedHash, hash) == 1
}
