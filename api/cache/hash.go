package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash calculates a short hex string for a given value.
// Optimized to avoid heap allocations by using unsafe for zero-allocation
// string-to-byte conversion and stack-allocating the hex string.
func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // Zero-allocation string to byte slice conversion for read-only hash
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
