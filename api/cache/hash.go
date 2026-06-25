package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	//nolint:gosec // G103: performance optimization for byte casting
	b := unsafe.Slice(unsafe.StringData(value), len(value))
	h := sha256.Sum256(b)

	// Optimize hex encoding to avoid allocations
	// 12 hex chars = 6 bytes from sha256
	var buf [12]byte
	hex.Encode(buf[:], h[:6])
	return string(buf[:])
}
