package cache

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func floatsToBytes(emb []float32) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, emb)
	return buf.Bytes()
}

func TestFloatsToBytes_RoundTrip(t *testing.T) {
	original := []float32{1.5, 2.5, -0.5, 1024.0}
	data := floatsToBytes(original)

	var restored []float32
	for i := 0; i < len(original); i++ {
		var v float32
		binary.Read(bytes.NewReader(data[i*4:]), binary.LittleEndian, &v)
		restored = append(restored, v)
	}

	if len(restored) != len(original) {
		t.Fatalf("len mismatch: got %d, want %d", len(restored), len(original))
	}
	for i := range restored {
		if restored[i] != original[i] {
			t.Errorf("mismatch at [%d]: got %v, want %v", i, restored[i], original[i])
		}
	}
}

func TestFloatsToBytes_512d(t *testing.T) {
	emb := make([]float32, 512)
	for i := range emb {
		emb[i] = float32(i) * 0.001
	}
	data := floatsToBytes(emb)
	if len(data) != 512*4 {
		t.Errorf("len=%d, want %d", len(data), 512*4)
	}
}
