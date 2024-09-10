package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/crypto/argon2"
)

// RandomPassword generates a random password of length 10
func RandomPassword() string {
	return RandomString(10)
}

// HashPasswordArgon2 hashes the password using Argon2id
func HashPasswordArgon2(password string) (string, error) {
	salt := make([]byte, 16) // 16 bytes salt
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	// Encode salt and hash for storage in a database
	saltStr := base64.RawStdEncoding.EncodeToString(salt)
	hashStr := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s:%s", saltStr, hashStr), nil
}

// Function to verify the password during login
func VerifyPassword(storedPassword string, providedPassword string) error {
	// Split the stored password into salt and hash
	parts := strings.Split(storedPassword, ":")
	if len(parts) != 2 {
		return errors.New("invalid stored password format")
	}
	saltStr := parts[0]
	storedHashStr := parts[1]

	// Decode the salt
	salt, err := base64.StdEncoding.DecodeString(saltStr)
	if err != nil {
		return fmt.Errorf("error decoding salt: %w", err)
	}

	// Hash the provided password using the same salt
	hashedPassword := hashPasswordWithSalt(providedPassword, salt)

	// Compare the new hash with the stored hash
	if hashedPassword != storedHashStr {
		return errors.New("password does not match")
	}

	// If they match, authentication is successful
	return nil
}

// Helper function to hash the provided password with the given salt (using Argon2)
func hashPasswordWithSalt(password string, salt []byte) string {
	// Argon2 parameters
	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	// Generate the Argon2 hash of the password with the provided salt
	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	// Encode the hash to a Base64 string
	hashStr := base64.StdEncoding.EncodeToString(hash)

	return hashStr
}
