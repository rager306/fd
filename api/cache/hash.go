package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash computes a SHA-256 hash of the input and returns the first 12 hex characters.
// It uses unsafe.Slice to avoid an allocation from casting the string to a byte slice,
// and it uses a stack-allocated buffer for hex encoding to avoid unnecessary heap allocations
// from hex.EncodeToString.
func shortHash(value string) string {
	var b []byte
	if len(value) > 0 {
		b = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(b)
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
