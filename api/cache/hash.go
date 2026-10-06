package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: Encode only the first 6 bytes to generate exactly 12 characters,
	// avoiding the allocation of a 64-byte string backing array.
	return hex.EncodeToString(h[:6])
}
