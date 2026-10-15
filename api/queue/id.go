package queue

import (
	"crypto/rand"
	"encoding/hex"
)

// NewRequestID returns a 16-hex-char identifier (64 bits of entropy)
// suitable for keying in-flight and completed queue items. Collisions
// are negligible in any fd deployment window. A failure of the system
// entropy source panics rather than silently returning a weak id.
func NewRequestID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic("queue: crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(buf[:])
}
