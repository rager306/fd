package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash calculates a 12-character hex representation of a SHA-256 hash.
// ⚡ Bolt: optimized to reduce heap allocations by using unsafe slice conversion
// and a stack-allocated byte array for hex encoding.
func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // safe usage for read-only hashing
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)

	// encode only the necessary bytes into a stack array to prevent heap allocation
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
