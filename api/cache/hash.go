package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	// Zero-copy string to byte slice conversion to avoid allocation.
	valueBytes := unsafe.Slice(unsafe.StringData(value), len(value)) //nolint:gosec // G103: performance optimization for byte casting
	h := sha256.Sum256(valueBytes)

	// Stack-allocate buffer for hex encoding to avoid hex.EncodeToString allocation.
	var buf [64]byte
	hex.Encode(buf[:], h[:])
	return string(buf[:12])
}
