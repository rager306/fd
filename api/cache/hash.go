package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Use a stack-allocated buffer to encode only the first 6 bytes (12 hex chars).
	// This avoids allocating a 64-byte string for the full hash before slicing.
	var buf [12]byte
	hex.Encode(buf[:], h[:6])
	return string(buf[:])
}
