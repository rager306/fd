package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash computes a shortened SHA256 hex digest of value.
// Optimized to reduce heap allocations by using unsafe slice cast
// and stack-allocated buffer for encoding.
func shortHash(value string) string {
	var bvalue []byte
	if len(value) > 0 {
		bvalue = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(bvalue)
	var dst [64]byte
	hex.Encode(dst[:], h[:])
	return string(dst[:12])
}
