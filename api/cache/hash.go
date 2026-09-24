package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // Safe zero-allocation string-to-byte conversion for read-only hashing
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)

	// Use a stack-allocated buffer for hex encoding to avoid heap allocation
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
