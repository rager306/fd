package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize allocation: encode required bytes into correctly sized stack-allocated array
	// before converting to string, preventing the entire 64-byte backing array from remaining in memory.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
