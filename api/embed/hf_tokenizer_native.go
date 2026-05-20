//go:build hf_tokenizers

package embed

import (
	"fmt"

	"github.com/daulet/tokenizers"
)

type nativeHFTokenizer struct {
	tokenizer *tokenizers.Tokenizer
}

type nativeHFEncoding struct {
	InputIDs      []int
	AttentionMask []int
}

func newNativeHFTokenizer(tokenizerPath string) (*nativeHFTokenizer, error) {
	if tokenizerPath == "" {
		return nil, fmt.Errorf("native hf tokenizer path is required")
	}
	tk, err := tokenizers.FromFile(tokenizerPath)
	if err != nil {
		return nil, fmt.Errorf("load native hf tokenizer %q: %w", tokenizerPath, err)
	}
	return &nativeHFTokenizer{tokenizer: tk}, nil
}

func (t *nativeHFTokenizer) Close() {
	if t == nil || t.tokenizer == nil {
		return
	}
	t.tokenizer.Close()
}

func (t *nativeHFTokenizer) Encode(text string) (nativeHFEncoding, error) {
	if t == nil || t.tokenizer == nil {
		return nativeHFEncoding{}, fmt.Errorf("native hf tokenizer is not initialized")
	}
	encoding := t.tokenizer.EncodeWithOptions(text, true, tokenizers.WithReturnAttentionMask())
	inputIDs := make([]int, len(encoding.IDs))
	for i, id := range encoding.IDs {
		inputIDs[i] = int(id)
	}
	attentionMask := make([]int, len(encoding.AttentionMask))
	for i, value := range encoding.AttentionMask {
		attentionMask[i] = int(value)
	}
	return nativeHFEncoding{InputIDs: inputIDs, AttentionMask: attentionMask}, nil
}
