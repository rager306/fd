package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: Optimize hex encoding to reduce heap allocations
	// and prevent memory leak of full 64-char string backing array
	var dst [64]byte
	hex.Encode(dst[:], h[:])
	return string(dst[:12])
}
