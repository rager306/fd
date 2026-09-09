package cache

import "testing"

func BenchmarkLocalCacheKey(b *testing.B) {
	key := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = localCacheKey(key, 1024)
	}
}

func BenchmarkLocalCacheKey_512(b *testing.B) {
	key := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = localCacheKey(key, 512)
	}
}

func BenchmarkLocalCacheKey_Other(b *testing.B) {
	key := "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = localCacheKey(key, 768)
	}
}
