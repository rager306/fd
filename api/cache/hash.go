package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: encode only the needed bytes into a stack-allocated buffer.
	// This avoids heap allocations from hex.EncodeToString and prevents
	// keeping the entire backing array alive in memory when truncating.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
