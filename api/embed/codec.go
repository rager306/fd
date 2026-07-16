// Package embed encodes and decodes embedding vectors for the OpenAI v1 encoding_format field (float array or base64 float32 LE bytes).
package embed

// Package embed encodes and decodes embedding vectors for the OpenAI v1 encoding_format field (float array or base64 float32 LE bytes).

import (
	"encoding/base64"
	"encoding/json"
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

func init() {
	// Verify system architecture is little endian.
	// Base64-encoded floats in the OpenAI format require little-endian.
	// The unsafe fast paths below assume the host architecture is also little-endian.
	buf := [2]byte{}
	//nolint:gosec // G103: fast path for safe slice conversion
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)
	if buf[0] == 0xAB {
		panic("embed codec: fast paths require little-endian architecture")
	}
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
	if len(slice) == 0 {
		return nil
	}
	byteLen := len(slice) * 4
	b := make([]byte, byteLen)

	// Create an aliased byte slice pointing to the float slice
	//nolint:gosec // G103: fast path for safe slice conversion, avoids aliasing by copying
	src := unsafe.Slice((*byte)(unsafe.Pointer(&slice[0])), byteLen)

	copy(b, src)
	return b
}

// BytesToFloat32Slice is the inverse of Float32SliceToBytes, used by tests
// and any future decode path (e.g. /v1/embeddings echo for symmetry).
func BytesToFloat32Slice(b []byte) []float32 {
	if len(b) == 0 || len(b)%4 != 0 {
		return nil
	}
	out := make([]float32, len(b)/4)

	// Create an aliased byte slice pointing to the float slice
	//nolint:gosec // G103: fast path for safe slice conversion, avoids aliasing by copying
	dst := unsafe.Slice((*byte)(unsafe.Pointer(&out[0])), len(b))

	copy(dst, b)
	return out
}
