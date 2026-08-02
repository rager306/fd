package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Optimize: allocate on stack and encode only needed prefix
	// to avoid hex.EncodeToString allocation and string sub-slicing leak.
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
