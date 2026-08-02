package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize hash string generation by encoding directly into a precisely
	// sized stack buffer to avoid keeping the full encoded backing array alive in memory.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
