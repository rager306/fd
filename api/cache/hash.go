package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// shortHash computes a fast sha256 hash of the value and returns a 12-char hex string.
// ⚡ Bolt: optimized by encoding only the first 6 bytes instead of encoding the full
// 32-byte hash (64-char string) and slicing it. This avoids retaining a larger 64-byte
// backing array in memory.
func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:6])
}
