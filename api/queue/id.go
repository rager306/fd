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
	// Optimization: Single-allocation string conversion for request IDs.
	var dst [16]byte
	hex.Encode(dst[:], buf[:])
	return string(dst[:])
}
