package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var h [32]byte
	if value == "" {
		h = sha256.Sum256(nil)
	} else {
		//nolint:gosec // Read-only conversion for hashing is safe.
		b := unsafe.Slice(unsafe.StringData(value), len(value))
		h = sha256.Sum256(b)
	}

	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
