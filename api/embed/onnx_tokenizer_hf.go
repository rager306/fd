//go:build hf_tokenizers

package embed

func newONNXTokenizer(tokenizerPath string) (onnxTokenizer, error) {
	return newNativeHFTokenizer(tokenizerPath)
}
