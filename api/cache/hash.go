package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	if value == "" {
		h := sha256.Sum256(nil)
		var buf [12]byte
		hex.Encode(buf[:], h[:6])
		return string(buf[:])
	}

	//nolint:gosec // G103: performance optimization for byte casting
	b := unsafe.Slice(unsafe.StringData(value), len(value))
	h := sha256.Sum256(b)

	var buf [12]byte
	hex.Encode(buf[:], h[:6])
	return string(buf[:])
}
