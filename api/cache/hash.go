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
	return hex.EncodeToString(h[:])[:12]
}
