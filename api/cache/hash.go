package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: Reduces allocations (from 2 to 1) and memory overhead
	// by stack-allocating the required 12 bytes instead of a full 64-character hex string.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
