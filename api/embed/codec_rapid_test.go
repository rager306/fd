package embed

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"testing"

	"pgregory.net/rapid"
)

func TestFloat32BytesRoundTrip_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		bits := rapid.SliceOfN(rapid.Uint32(), 0, 128).Draw(t, "bits")
		values := make([]float32, len(bits))
		for i, b := range bits {
			values[i] = math.Float32frombits(b)
		}

		encoded := Float32SliceToBytes(values)
		if got, want := len(encoded), len(values)*4; got != want {
			t.Fatalf("encoded byte length = %d, want %d", got, want)
		}

		decoded := BytesToFloat32Slice(encoded)
		if got, want := len(decoded), len(values); got != want {
			t.Fatalf("decoded length = %d, want %d", got, want)
		}
		for i := range values {
			if got, want := math.Float32bits(decoded[i]), math.Float32bits(values[i]); got != want {
				t.Fatalf("decoded[%d] bits = %08x, want %08x", i, got, want)
			}
		}
	})
}

func TestBytesToFloat32SliceRejectsMisalignedInput_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		base := rapid.SliceOfN(rapid.Byte(), 0, 64).Draw(t, "base")
		extra := rapid.IntRange(1, 3).Draw(t, "extra")
		misaligned := append(append([]byte(nil), base...), make([]byte, extra)...)
		if len(misaligned)%4 == 0 {
			misaligned = append(misaligned, 0)
		}

		if decoded := BytesToFloat32Slice(misaligned); decoded != nil {
			t.Fatalf("BytesToFloat32Slice(%d misaligned bytes) = %#v, want nil", len(misaligned), decoded)
		}
	})
}

func TestEncodeEmbeddingBase64MatchesByteCodec_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		values := rapid.SliceOfN(rapid.Float32Range(-1_000_000, 1_000_000), 0, 128).Draw(t, "values")

		encoded := EncodeEmbedding(values, EncodingFormatBase64)
		decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			t.Fatalf("base64 decode failed: %v", err)
		}
		decoded := BytesToFloat32Slice(decodedBytes)
		if got, want := len(decoded), len(values); got != want {
			t.Fatalf("decoded length = %d, want %d", got, want)
		}
		for i := range values {
			if decoded[i] != values[i] {
				t.Fatalf("decoded[%d] = %v, want %v", i, decoded[i], values[i])
			}
		}
	})
}

func TestEncodeEmbeddingDefaultJSONRoundTrip_Rapid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		values := rapid.SliceOfN(rapid.Float32Range(-1_000_000, 1_000_000), 0, 128).Draw(t, "values")

		encoded := EncodeEmbedding(values, "")
		var decoded []float32
		if err := json.Unmarshal([]byte(encoded), &decoded); err != nil {
			t.Fatalf("json decode failed: %v", err)
		}
		if got, want := len(decoded), len(values); got != want {
			t.Fatalf("decoded length = %d, want %d", got, want)
		}
		for i := range values {
			if decoded[i] != values[i] {
				t.Fatalf("decoded[%d] = %v, want %v", i, decoded[i], values[i])
			}
		}
	})
}
