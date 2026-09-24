package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: Optimize shortHash to prevent heap allocations
	// By encoding directly into a stack-allocated byte array instead of
	// using hex.EncodeToString() and slicing, we avoid allocating a full 64-byte
	// string on the heap only to throw most of it away.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
