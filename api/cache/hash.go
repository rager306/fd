package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash computes a 12-character hex string from the sha256 hash of value.
// It is optimized to minimize allocations on the hot path (cache lookups).
// Performance impact: Reduces allocations from 3 to 1 per call by using
// unsafe slice conversion and a stack-allocated buffer for hex encoding.
func shortHash(value string) string {
	var b []byte
	//nolint:gocritic // intentionally checking length for unsafe slice creation
	if len(value) > 0 {
		//nolint:gosec // zero-allocation string to byte slice conversion for read-only hashing in hot path
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)

	// Stack-allocate the destination buffer for the first 6 bytes (12 hex chars).
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
