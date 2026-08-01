package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: Avoid allocating the full 64-character hex string which would be kept
	// alive in memory when we slice it. Instead, encode directly to a fixed-size stack buffer.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
