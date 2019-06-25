package uuid

import (
	"github.com/btcsuite/btcutil/base58"
	google_uuid "github.com/google/uuid"
)

func NewBase58() (string, error) {
	id, err := google_uuid.NewUUID()
	if err != nil {
		return "", err
	}
	b, err := id.MarshalBinary()
	if err != nil {
		return "", err
	}

	s := base58.Encode(b)

	return s, err
}
