package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash calculates a 12-character hex hash.
// Uses unsafe string conversion and stack allocation to minimize heap escapes.
func shortHash(value string) string {
	var b []byte
	if value != "" {
		//nolint:gosec // Safe zero-allocation conversion for read-only hash operation
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
