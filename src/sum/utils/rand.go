package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func RandBytes(length int) []byte {
	buf := make([]byte, length)
	rand.Reader.Read(buf)

	return buf
}

func RandBase64(bytes int) string {
	return base64.RawURLEncoding.EncodeToString(RandBytes(12))
}
