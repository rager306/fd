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

func BenchmarkRedisCacheKey(b *testing.B) {
	c := &RedisCache{prefix: "bench:", namespace: "v2"}
	text := "benchmark test string for key performance"

	b.Run("dim1024", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = c.key(text, 1024)
		}
	})

	b.Run("dim512", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = c.key(text, 512)
		}
	})

	b.Run("dim256", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = c.key(text, 256)
		}
	})
}

func BenchmarkHashText_Short(b *testing.B) {
	c := &RedisCache{prefix: "bench:"}
	text := "hi"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.HashText(text)
	}
}
