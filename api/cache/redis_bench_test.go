package cache

import (
	"testing"
)

func BenchmarkHashText(b *testing.B) {
	c := &RedisCache{}
	text := "this is a typical string that we want to hash for our cache key generation to measure performance."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.HashText(text)
	}
}

func BenchmarkRedisKey(b *testing.B) {
	c := &RedisCache{prefix: "embed:cache:", namespace: "v2:m1234567:r1234567"}
	text := "this is a typical string that we want to hash for our cache key generation to measure performance."
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.key(text, 1024)
	}
}
