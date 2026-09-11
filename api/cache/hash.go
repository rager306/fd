package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash generates a stable, short 12-character hex hash from the given string.
// Optimized to reduce heap allocations on hot paths.
//
//nolint:gosec // unsafe string to byte slice conversion is safe for read-only hashing
func shortHash(value string) string {
	var b []byte
	if value != "" {
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
