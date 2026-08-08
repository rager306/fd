package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash calculates a SHA256 hash of the input string and returns
// the first 12 characters of its hex encoding.
func shortHash(value string) string {
	var b []byte
	if value != "" {
		// Zero-allocation string-to-byte conversion for read-only hash operation.
		//nolint:gosec // Memory aliasing is safe here as sha256.Sum256 only reads the bytes
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)

	// Stack-allocate to avoid hex string and sub-slice allocations
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
