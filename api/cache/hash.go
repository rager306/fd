package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// ⚡ Bolt: encode only the first 6 bytes to directly get a 12-char hex string
	// instead of encoding 32 bytes and slicing the result, preventing memory bloat.
	return hex.EncodeToString(h[:6])
}
