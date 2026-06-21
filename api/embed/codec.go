// Package embed encodes and decodes embedding vectors for the OpenAI v1 encoding_format field (float array or base64 float32 LE bytes).
package embed

// Package embed encodes and decodes embedding vectors for the OpenAI v1 encoding_format field (float array or base64 float32 LE bytes).

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math"
	"unsafe"
)

// Encoding format constants. Used by /v1/embeddings and /embeddings/batch
// to choose between float arrays and base64-encoded float32 LE bytes.
// Base64 saves ~30% bandwidth (4 bytes float32 → ~5.4 chars base64 per
// element vs ~14 chars for JSON-encoded float).
const (
	EncodingFormatFloat  = "float"
	EncodingFormatBase64 = "base64"
)

var isLittleEndian bool

func init() {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i) //nolint:gosec // G103: performance optimization for byte casting
	pb := (*byte)(u)
	isLittleEndian = (*pb == 0x04)
}

// EncodeEmbedding serializes an embedding vector in the requested format.
// `format` is one of EncodingFormatFloat or EncodingFormatBase64; the empty
// string defaults to float. Any other value returns the float form (callers
// should validate format before calling).
func EncodeEmbedding(emb []float32, format string) string {
	if format == EncodingFormatBase64 {
		return base64.StdEncoding.EncodeToString(Float32SliceToBytes(emb))
	}
	b, _ := json.Marshal(emb)
	return string(b)
}

// Float32SliceToBytes converts a float32 slice to a little-endian byte
// slice suitable for base64 encoding. Length must equal len(slice)*4.
func Float32SliceToBytes(slice []float32) []byte {
	if isLittleEndian && len(slice) > 0 {
		bytesLen := len(slice) * 4
		out := make([]byte, bytesLen)
		src := unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), bytesLen) //nolint:gosec // G103: performance optimization for byte casting
		copy(out, src)
		return out
	}

	b := make([]byte, len(slice)*4)
	for i, v := range slice {
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(v))
	}
	return b
}

// BytesToFloat32Slice is the inverse of Float32SliceToBytes, used by tests
// and any future decode path (e.g. /v1/embeddings echo for symmetry).
func BytesToFloat32Slice(b []byte) []float32 {
	if len(b)%4 != 0 {
		return nil
	}
	out := make([]float32, len(b)/4)
	for i := range out {
		bits := binary.LittleEndian.Uint32(b[i*4:])
		out[i] = math.Float32frombits(bits)
	}
	return out
}
