package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: encode only the required 6 bytes into a stack-allocated array
	// to avoid allocating a full 64-byte string on the heap before slicing.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
