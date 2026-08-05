package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimization: Single-allocation string conversion for shortened hashes.
	// Only encode the required prefix to avoid memory leaks from string slicing.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
