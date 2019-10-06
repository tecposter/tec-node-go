package uuid

import (
	"github.com/btcsuite/btcutil/base58"
	google_uuid "github.com/google/uuid"
	"github.com/tecposter/tec-node-go/lib/dto"
)

// NewBase58 returns a uuid as base58 encoded string
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

// NewID returns a uuid wrapped by dto.ID
func NewID() (dto.ID, error) {
	id, err := google_uuid.NewUUID()
	if err != nil {
		return dto.ID([]byte{}), err
	}
	b, err := id.MarshalBinary()
	return dto.ID(b), err
}
