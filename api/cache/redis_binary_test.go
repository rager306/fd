package cache

import (
	"testing"
)

func TestBinaryMarshalUnmarshal_4d(t *testing.T) {
	emb := []float32{1.0, 2.0, 3.0, 4.0}
	data := marshalEmbedding(emb, 4)

	if len(data) != 2+4*4 {
		t.Errorf("len=%d, want %d", len(data), 2+16)
	}

	got, dim := unmarshalEmbedding(data)
	if dim != 4 {
		t.Errorf("dim=%d, want 4", dim)
	}
	if len(got) != 4 {
		t.Errorf("len=%d, want 4", len(got))
	}
	for i := range got {
		if got[i] != emb[i] {
			t.Errorf("got[%d]=%v, want %v", i, got[i], emb[i])
		}
	}
}

func TestBinaryMarshalUnmarshal_512d(t *testing.T) {
	emb := make([]float32, 512)
	for i := range emb {
		emb[i] = float32(i) * 0.001
	}
	data := marshalEmbedding(emb, 512)
	got, dim := unmarshalEmbedding(data)

	if dim != 512 {
		t.Errorf("dim=%d, want 512", dim)
	}
	if len(got) != 512 {
		t.Errorf("len=%d, want 512", len(got))
	}
}

func TestBinaryMarshalUnmarshal_1024d(t *testing.T) {
	emb := make([]float32, 1024)
	for i := range emb {
		emb[i] = float32(i) * 0.001
	}
	data := marshalEmbedding(emb, 1024)
	if len(data) != 2+1024*4 {
		t.Errorf("len=%d, want %d", len(data), 2+4096)
	}

	got, dim := unmarshalEmbedding(data)
	if dim != 1024 {
		t.Errorf("dim=%d, want 1024", dim)
	}
	if len(got) != 1024 {
		t.Errorf("len=%d, want 1024", len(got))
	}
}

func TestBinaryMarshalUnmarshal_Truncated(t *testing.T) {
	_, dim := unmarshalEmbedding([]byte{1}) // too short
	if dim != 0 {
		t.Errorf("dim=%d, want 0 for truncated", dim)
	}
}

func TestBinaryMarshalUnmarshal_RoundTrip(t *testing.T) {
	original := []float32{0.123456, -0.654321, 1e10, -1e-10}
	data := marshalEmbedding(original, 4)
	got, _ := unmarshalEmbedding(data)

	if len(got) != len(original) {
		t.Fatalf("len mismatch: got %d, want %d", len(got), len(original))
	}
	for i := range got {
		if got[i] != original[i] {
			t.Errorf("roundtrip failed at [%d]: got %v, want %v", i, got[i], original[i])
		}
	}
}
