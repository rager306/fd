package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	if value == "" {
		h := sha256.Sum256(nil)
		var buf [64]byte
		hex.Encode(buf[:], h[:])
		return string(buf[:12])
	}
	//nolint:gosec // G103: performance optimization for byte casting
	b := unsafe.Slice(unsafe.StringData(value), len(value))
	h := sha256.Sum256(b)
	var buf [64]byte
	hex.Encode(buf[:], h[:])
	return string(buf[:12])
}
