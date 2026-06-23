package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var valBytes []byte
	if value != "" {
		valBytes = unsafe.Slice(unsafe.StringData(value), len(value)) //nolint:gosec // G103: performance optimization for byte casting
	}
	h := sha256.Sum256(valBytes)
	var buf [12]byte
	hex.Encode(buf[:], h[:6])
	return string(buf[:])
}
