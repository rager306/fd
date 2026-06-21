package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	//nolint:gosec // G103: performance optimization for byte casting
	h := sha256.Sum256(unsafe.Slice(unsafe.StringData(value), len(value)))
	return hex.EncodeToString(h[:])[:12]
}
