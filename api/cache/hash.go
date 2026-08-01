package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: allocate exact-size stack buffer to avoid memory leak
	// from subslicing the full hex string backing array and reduce heap allocations.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
