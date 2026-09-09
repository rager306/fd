package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	if value == "" {
		return "e3b0c44298fc" // precomputed empty hash prefix
	}
	// Zero-allocation string to byte slice
	//nolint:gosec // Safe because the slice is only read by sha256.Sum256
	b := unsafe.Slice(unsafe.StringData(value), len(value))
	h := sha256.Sum256(b)

	// Encode directly into a stack-allocated array to prevent allocating a full hex string
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
