package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize memory allocation by only hex encoding the first 6 bytes (12 hex chars)
	// instead of encoding the full 32-byte hash and slicing the string, which would
	// retain the entire 64-byte string backing array in memory.
	return hex.EncodeToString(h[:6])
}
