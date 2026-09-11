package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash generates a 12-char hex string hash.
// Optimized to avoid heap allocations for string conversion and hex encoding.
func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // zero-allocation string to byte slice conversion for read-only hashing
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
