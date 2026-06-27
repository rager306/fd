package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var bytes []byte
	if value != "" {
		bytes = unsafe.Slice(unsafe.StringData(value), len(value)) //nolint:gosec // G103: performance optimization for byte casting
	}
	h := sha256.Sum256(bytes)
	var buf [64]byte
	hex.Encode(buf[:], h[:])
	return string(buf[:12])
}
