package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	var dst [64]byte
	hex.Encode(dst[:], h[:])
	return string(dst[:12])
}
