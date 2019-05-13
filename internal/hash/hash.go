package hash

import (
	"crypto/sha256"
	b58 "github.com/mr-tron/base58/base58"
)

func GetContentId(content string) string {
	h := sha256.New()
	h.Write([]byte(content))

	return b58.Encode(h.Sum(nil))
}
