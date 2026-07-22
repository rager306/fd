package embed

import (
	"testing"
)

func BenchmarkFloat32SliceToBytes(b *testing.B) {
	slice := make([]float32, 1024)
	for i := range slice {
		slice[i] = float32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Float32SliceToBytes(slice)
	}
}

func BenchmarkBytesToFloat32Slice(b *testing.B) {
	slice := make([]float32, 1024)
	for i := range slice {
		slice[i] = float32(i)
	}
	data := Float32SliceToBytes(slice)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BytesToFloat32Slice(data)
	}
}
