package rand

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateBytes(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateStr(size int) (string, error) {
	b, err := GenerateBytes(size)
	if err != nil {
		return "", nil
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
