package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	// Encode only the needed bytes to avoid memory bloat from slicing a large string.
	return hex.EncodeToString(h[:6])
}
