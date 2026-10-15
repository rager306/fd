package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

// shortHash returns the first 12 hex characters of the SHA-256 digest of
// value. The 12-char length is part of the cache-key and Redis namespace
// format, so it must stay stable.
func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:])
}
