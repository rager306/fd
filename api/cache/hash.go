package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash returns a 12-character hex representation of the SHA256 hash of the value.
// Optimized to avoid heap allocations by using unsafe for zero-allocation string to byte slice conversion,
// and hex encoding directly into a stack-allocated byte array, rather than creating a 64-byte string
// dynamically and keeping the entire backing array in memory.
func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // Zero-allocation conversion for hashing
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
