package cache

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(value string) string {
	h := sha256.Sum256([]byte(value))
	var buf [12]byte
	hex.Encode(buf[:], h[:6])
	return string(buf[:])
}
