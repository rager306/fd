package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"unsafe"
)

func BenchmarkHashTextOriginal(b *testing.B) {
	text := "example text to hash for the embedding cache which could be long"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h := sha256.Sum256([]byte(text))
		_ = hex.EncodeToString(h[:])
	}
}

func BenchmarkHashTextOptimized(b *testing.B) {
	text := "example text to hash for the embedding cache which could be long"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var textBytes []byte
		if text != "" {
			textBytes = unsafe.Slice(unsafe.StringData(text), len(text))
		}
		h := sha256.Sum256(textBytes)
		var buf [64]byte
		hex.Encode(buf[:], h[:])
		_ = string(buf[:])
	}
}
