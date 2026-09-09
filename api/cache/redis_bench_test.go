package cache

import (
	"testing"
)

func BenchmarkHashText(b *testing.B) {
	c := &RedisCache{prefix: "bench:"}
	text := "benchmark test string for hashing performance"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.HashText(text)
	}
}

func BenchmarkHashText_Short(b *testing.B) {
	c := &RedisCache{prefix: "bench:"}
	text := "hi"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.HashText(text)
	}
}

func BenchmarkRedisKey(b *testing.B) {
	c := &RedisCache{
		prefix:    "embed:cache:",
		namespace: "v2",
	}
	text := "This is a sample text for hashing and generating a key."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.key(text, 1024)
	}
}
