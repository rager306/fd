package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

// shortHash computes a SHA256 hash and returns the first 12 hex characters.
// It is optimized to avoid allocating a full 64-byte hex string and leaking memory
// when sliced.
func shortHash(value string) string {
	var h [32]byte
	if value != "" {
		h = sha256.Sum256(unsafe.Slice(unsafe.StringData(value), len(value))) //nolint:gosec // use of unsafe is audited
	} else {
		h = sha256.Sum256(nil)
	}
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
