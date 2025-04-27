package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/argon2"
)

const (
	hashMemory      = 1024 * 32 // 32 MiB
	hashIterations  = 1
	hashParallelism = 2
	saltLength      = 8
	keyLength       = 32
)

// NewSalt generates a random salt of length saltLength, the result is a byte slice
func NewSalt() ([]byte, error) {
	return generateRandomBytes(saltLength)
}

// HashPassword hashes a password using Argon2 with a given salt, the result is a base64 representation of the hashed pa
func HashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, hashIterations, hashMemory, hashParallelism, keyLength)
}

func Compare(password string, base64Hash string, salt []byte) (bool, error) {
	hash, err := base64.RawStdEncoding.Strict().DecodeString(base64Hash)
	if err != nil {
		return false, err
	}

	newHash := HashPassword(password, salt)
	return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func Check(password, hashedPassword string) (bool, error) {
	// Decode the hashed password
	decoded, err := base64.RawStdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false, err
	}

	// Extract the salt (first 16 bytes)
	if len(decoded) < saltLength {
		return false, errors.New("invalid hashed password")
	}
	salt := decoded[:saltLength]
	hashed := decoded[saltLength:]

	newHash := HashPassword(password, salt)
	return subtle.ConstantTimeCompare(hashed, newHash) == 1, nil
}
