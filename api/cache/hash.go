package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // zero-allocation string-to-byte conversion for read-only hashing in hot path
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)

	// Optimize hash string creation to reduce heap allocations.
	// Encode only the first 6 bytes into a stack-allocated array
	// instead of allocating a full 64-char string on the heap.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
