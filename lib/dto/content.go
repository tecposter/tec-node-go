package dto

import (
	"crypto/sha256"
)

// GenContentID returns dto.ID from content using sha256
func GenContentID(content string) ID {
	h := sha256.New()
	h.Write([]byte(content))
	return ID(h.Sum(nil))
}
