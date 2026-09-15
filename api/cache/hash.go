package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash computes a 12-character hex prefix of the SHA-256 hash of the value.
func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // Read-only conversion for zero-allocation performance in hot path
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)

	// Avoid hex.EncodeToString heap allocation and string slice allocation
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
