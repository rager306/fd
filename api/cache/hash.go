package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: Zero-allocation fast path for short hash.
	// Encodes directly into a stack-allocated array instead of dynamically
	// allocating a 64-byte string on the heap only to slice the first 12 bytes.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
