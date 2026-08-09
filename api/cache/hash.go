package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"
)

func shortHash(value string) string {
	var b []byte //nolint:gosec // Memory aliasing is safe for read-only hashing
	if len(value) > 0 { //nolint:gocritic // emptyStringTest does not apply as len is needed
		b = unsafe.Slice(unsafe.StringData(value), len(value)) //nolint:gosec // Memory aliasing is safe for read-only hashing
	}
	h := sha256.Sum256(b)
	var dst [12]byte
	hex.Encode(dst[:], h[:6])
	return string(dst[:]) //nolint:gosec // Memory aliasing is safe for read-only hashing
}
