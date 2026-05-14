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