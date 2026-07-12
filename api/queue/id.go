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
	_, _ = rand.Read(buf[:])
	var hexBuf [16]byte; hex.Encode(hexBuf[:], buf[:]); return string(hexBuf[:])
}
