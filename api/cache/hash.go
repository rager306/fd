package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize short hash generation by avoiding heap allocations
	// for the full 64-character hex string and slicing.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
