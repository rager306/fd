package cache

import "testing"

func BenchmarkShortHashOriginal(b *testing.B) {
	val := "some_cache_key_string"
	for i := 0; i < b.N; i++ {
		_ = shortHashOriginal(val)
	}
}

func BenchmarkShortHashOptimized(b *testing.B) {
	val := "some_cache_key_string"
	for i := 0; i < b.N; i++ {
		_ = shortHash(val)
	}
}
