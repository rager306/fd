package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: use stack-allocated buffer to encode only required prefix, avoiding memory leak from sub-slicing
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
