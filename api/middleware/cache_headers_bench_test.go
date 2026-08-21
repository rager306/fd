package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func responseETagOld(body []byte) string {
	sum := sha256.Sum256(body)
	return `"` + hex.EncodeToString(sum[:]) + `"`
}

func responseETagNew(body []byte) string {
	sum := sha256.Sum256(body)
	var dst [66]byte
	dst[0] = '"'
	hex.Encode(dst[1:65], sum[:])
	dst[65] = '"'
	return string(dst[:])
}

func BenchmarkResponseETagOld(b *testing.B) {
	body := []byte("test body data for etag benchmark")
	for i := 0; i < b.N; i++ {
		_ = responseETagOld(body)
	}
}

func BenchmarkResponseETagNew(b *testing.B) {
	body := []byte("test body data for etag benchmark")
	for i := 0; i < b.N; i++ {
		_ = responseETagNew(body)
	}
}
