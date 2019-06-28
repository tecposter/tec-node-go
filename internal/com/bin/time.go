package bin

import (
	"time"
)

// BytesToTime converts a byte slide to Time
func BytesToTime(buf []byte) time.Time {
	nsec := BytesToInt64(buf)
	return time.Unix(0, nsec)
}

// TimeToBytes converts Time to a byte slide with size 8
func TimeToBytes(t time.Time) []byte {
	nano := t.UnixNano()
	return Int64ToBytes(nano)
}
