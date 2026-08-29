package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var input []byte
	if value != "" {
		//nolint:gosec // Safe read-only conversion for hashing
		input = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(input)

	// Optimize to reduce allocations:
	// Encode directly to a stack-allocated array instead of a 64-char string.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
