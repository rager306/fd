package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize short hash generation by only hex-encoding the required bytes
	// into a stack-allocated buffer to eliminate a heap allocation and avoid
	// the Go string slicing memory leak that keeps the full backing array alive.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
