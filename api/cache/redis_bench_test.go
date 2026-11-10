package cache

import "testing"

var benchmarkHashTextResult string

func BenchmarkHashText(b *testing.B) {
	c := &RedisCache{prefix: "bench:"}
	text := "benchmark test string for hashing performance"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkHashTextResult = c.HashText(text)
	}
}

func BenchmarkHashText_Short(b *testing.B) {
	c := &RedisCache{prefix: "bench:"}
	text := "hi"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchmarkHashTextResult = c.HashText(text)
	}
}
