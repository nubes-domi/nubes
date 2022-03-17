package utils

import (
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	seed := RandBytes(32)
	hash := HashPasswordWithSeed(password, seed)
	return base64.RawStdEncoding.EncodeToString(hash) + "$" + base64.RawStdEncoding.EncodeToString(seed)
}

func HashPasswordWithSeed(password string, seed []byte) []byte {
	return argon2.IDKey([]byte(password), seed, 1, 64*1024, 4, 32)
}
