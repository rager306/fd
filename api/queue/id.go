package queue

import (
	"crypto/rand"
	"encoding/hex"
)

// NewRequestID returns a 16-hex-char identifier (64 bits of entropy)
// suitable for keying in-flight and completed queue items. Collisions
// are negligible in any fd deployment window.
func NewRequestID() string {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		panic("crypto/rand failed to provide entropy: " + err.Error())
	}
	return hex.EncodeToString(buf[:])
}
