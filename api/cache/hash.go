package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var valueBytes []byte
	if value != "" {
		//nolint:gosec // G103: performance optimization for byte casting
		valueBytes = unsafe.Slice(unsafe.StringData(value), len(value))
	}
	h := sha256.Sum256(valueBytes)
	var buf [64]byte
	hex.Encode(buf[:], h[:])
	return string(buf[:12])
}
