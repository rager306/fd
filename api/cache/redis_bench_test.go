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

func BenchmarkMarshalEmbedding(b *testing.B) {
	slice := make([]float32, 1024)
	for i := range slice {
		slice[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marshalEmbedding(slice, 1024)
	}
}

func BenchmarkUnmarshalEmbedding(b *testing.B) {
	slice := make([]float32, 1024)
	for i := range slice {
		slice[i] = float32(i)
	}
	data, _ := marshalEmbedding(slice, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		unmarshalEmbedding(data)
	}
}
