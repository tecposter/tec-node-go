package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateBytes generates random bytes with specific size
func GenerateBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateStr generates random string with specific byte size
func GenerateStr(size int) (string, error) {
	b, err := GenerateBytes(size)
	if err != nil {
		return "", nil
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
