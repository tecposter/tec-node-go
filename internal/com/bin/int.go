package bin

import (
	"encoding/binary"
)

// BytesToInt64 converts byte slide to int64
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// Int64ToBytes converts int64 to byte slide
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
