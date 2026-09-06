package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: Avoid allocating a full 64-char string on the heap.
	// Encode only the 6 required bytes into a small stack-allocated array.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
