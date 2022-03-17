package utils

import "crypto/rand"

func RandBytes(length int) []byte {
	buf := make([]byte, length)
	rand.Reader.Read(buf)

	return buf
}
