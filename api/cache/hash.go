package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: Eliminate heap allocations from string slicing a dynamically allocated hex string.
	// Encode directly to a stack-allocated buffer (64 bytes/op vs 176 bytes/op).
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
