//go:build onnx || hf_tokenizers

package embed

type onnxTokenEncoding struct {
	InputIDs      []int
	AttentionMask []int
}

type onnxTokenizer interface {
	Encode(text string) (onnxTokenEncoding, error)
	Close()
}
