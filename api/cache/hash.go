package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize hot path: hex.EncodeToString allocates on heap and slicing it
	// leaks the entire 64-byte backing array. Use a stack buffer and encode
	// only the necessary 6 bytes (12 hex chars).
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
