package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// shortHash generates a 12-character hex string hash.
// It encodes directly into a stack-allocated buffer and casts to string
// to avoid heap allocations and memory leaks from slicing a larger string.
func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
