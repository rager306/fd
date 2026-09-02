package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash creates a 12-character hex representation of a SHA256 hash.
// Optimized to avoid unnecessary heap allocations during string/byte conversion and hex encoding.
func shortHash(value string) string {
	var input []byte
	if value != "" {
		//nolint:gosec // Zero-allocation string-to-byte conversion for read-only hashing.
		input = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(input)

	// Encode only the required 6 bytes into a stack-allocated array instead of the full 32 bytes (64 chars).
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
