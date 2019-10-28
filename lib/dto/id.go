package dto

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/btcsuite/btcutil/base58"
)

// A ID contains functions to represent data in different formats
type ID []byte

// EmptyID returns empty ID
func EmptyID() ID {
	return ID([]byte{})
}

// Base58ToID returns ID converted from base58 encoding
func Base58ToID(src string) ID {
	if src == "" {
		return nil
	}
	d := base58.Decode(string(src))
	return ID(d)
}

// Bytes returns byte slide of a id
func (id ID) Bytes() []byte {
	return []byte(id)
}

// Base58 returns the base58 encoding of a id
func (id ID) Base58() string {
	return base58.Encode(id)
}

// Base64 returns the base64 encoding of a id
func (id ID) Base64() string {
	return base64.StdEncoding.EncodeToString(id)
}

// MarshalJSON implements the json.Marshaler interface
func (id ID) MarshalJSON() ([]byte, error) {
	jsonVal, err := json.Marshal(id.Base58())
	return jsonVal, err
}

// Equal reports whether id and b are the same length and contain the same bytes. A nil argument is
// equivalent to an empty slice.
func (id ID) Equal(b ID) bool {
	return bytes.Equal(id, b)
}

// UnmarshalJSON implements the json.UnmarshalJSON interface
/*
func (id *ID) UnmarshalJSON(b []byte) error {
	d := base58.Decode(string(b))
	*id = ID(d)
	return nil
}
*/
