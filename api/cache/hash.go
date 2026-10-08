package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// shortHash computes a SHA256 hash and returns the first 12 hex characters.
// Slicing the hash before hex encoding reduces memory allocations and avoids
// retaining the full 64-character hex string backing array.
func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	return hex.EncodeToString(h[:6])
}
