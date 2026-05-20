//go:build onnx && !hf_tokenizers

package embed

import (
	"fmt"

	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

type sugarmeONNXTokenizer struct {
	tokenizer *tokenizer.Tokenizer
}

func newONNXTokenizer(tokenizerPath string) (onnxTokenizer, error) {
	tk, err := pretrained.FromFile(tokenizerPath)
	if err != nil {
		return nil, fmt.Errorf("load onnx tokenizer %q: %w", tokenizerPath, err)
	}
	return &sugarmeONNXTokenizer{tokenizer: tk}, nil
}

func (t *sugarmeONNXTokenizer) Encode(text string) (onnxTokenEncoding, error) {
	encoding, err := t.tokenizer.EncodeSingle(text, true)
	if err != nil {
		return onnxTokenEncoding{}, err
	}
	inputIDs := make([]int, len(encoding.Ids))
	copy(inputIDs, encoding.Ids)
	attentionMask := make([]int, len(encoding.AttentionMask))
	copy(attentionMask, encoding.AttentionMask)
	return onnxTokenEncoding{InputIDs: inputIDs, AttentionMask: attentionMask}, nil
}

func (t *sugarmeONNXTokenizer) Close() {}
