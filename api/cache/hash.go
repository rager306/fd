package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: Avoid allocating the full 64-character hex string (128 bytes) and slicing it.
	// We only need 12 hex characters, which requires 6 bytes from the hash.
	// This reduces allocations from 2 to 1 and prevents the large backing array from being retained in memory.
	return hex.EncodeToString(h[:6])
}
