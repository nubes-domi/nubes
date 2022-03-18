package utils

import (
	"bytes"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	seed := RandBytes(32)
	hash := hashPasswordWithSeed(password, seed)
	return base64.RawStdEncoding.EncodeToString(hash) + "$" + base64.RawStdEncoding.EncodeToString(seed)
}

func VerifyPassword(password string, digest string) bool {
	pieces := strings.Split(digest, "$")

	seed, _ := base64.RawStdEncoding.DecodeString(pieces[1])
	expected, _ := base64.RawStdEncoding.DecodeString(pieces[0])
	actual := hashPasswordWithSeed(password, seed)

	return bytes.Equal(actual, expected)
}

func hashPasswordWithSeed(password string, seed []byte) []byte {
	return argon2.IDKey([]byte(password), seed, 1, 64*1024, 4, 32)
}
