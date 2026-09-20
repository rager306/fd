package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	if value == "" {
		h := sha256.Sum256([]byte{})
		var dst [12]byte
		hex.Encode(dst[:], h[:6])
		return string(dst[:])
	}

	// ⚡ Bolt: Zero-allocation string to byte slice conversion for read-only hashing.
	//nolint:gosec // Safe because the byte slice is only used for read-only hashing and does not escape.
	b := unsafe.Slice(unsafe.StringData(value), len(value))
	h := sha256.Sum256(b)

	// ⚡ Bolt: Encode directly to stack-allocated array instead of hex.EncodeToString to avoid heap allocation.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
