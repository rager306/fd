package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	bSlice := unsafe.Slice(unsafe.StringData(value), len(value)) //nolint:gosec // G103: performance optimization for byte casting
	h := sha256.Sum256(bSlice)
	var buf [12]byte
	hex.Encode(buf[:], h[:6])
	return string(buf[:])
}
