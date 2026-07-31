package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: encode only the required prefix into a stack-allocated buffer
	// to avoid heap allocation and prevent keeping the entire 64-byte underlying array alive.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
