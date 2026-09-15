package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// shortHash returns a 12-character hex string of the SHA-256 hash.
// Optimized to encode only the required bytes into a stack-allocated array,
// reducing allocations and preventing memory bloat from string slicing.
func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
