package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: encode only required prefix into stack-allocated buffer
	// to avoid heap allocation and backing array memory leak.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
