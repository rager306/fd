package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// shortHash creates a 12-character hex encoded string from the SHA256 of the value.
// It uses a stack allocated byte array to avoid unnecessary heap allocations during hex encoding.
func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
