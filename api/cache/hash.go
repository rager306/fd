package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// shortHash generates a 12-character hex hash.
// ⚡ Bolt: Optimized to use a stack-allocated buffer and avoid keeping the full 64-byte hex string in memory.
func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
