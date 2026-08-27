package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash generates a 12-character hex hash.
// Optimized: Reduces allocations by using zero-allocation string-to-byte conversion
// and encoding directly into a stack-allocated byte array (3 allocs -> 1 alloc).
func shortHash(value string) string {
	var h [32]byte
	if len(value) > 0 { //nolint:gocritic // needed for safety with unsafe string slice
		h = sha256.Sum256(unsafe.Slice(unsafe.StringData(value), len(value))) //nolint:gosec // zero-allocation string-to-byte-slice for read-only hash
	} else {
		h = sha256.Sum256(nil)
	}
	var dst [12]byte
	hex.Encode(dst[:], h[:6]) // Only encode the first 6 bytes to get 12 hex chars
	return string(dst[:])
}
