package middleware

import "testing"

func BenchmarkResponseETag(b *testing.B) {
	val := []byte("this is a typical string for hashing that we want to benchmark and optimize for performance")
	for i := 0; i < b.N; i++ {
		_ = responseETag(val)
	}
}
