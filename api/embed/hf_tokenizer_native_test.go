//go:build hf_tokenizers

package embed

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type tokenBaselineProbe struct {
	Label         string `json:"label"`
	InputIDs      []int  `json:"input_ids"`
	AttentionMask []int  `json:"attention_mask"`
}

func TestNativeHFTokenizerMatchesBaseline(t *testing.T) {
	tokenizerPath := os.Getenv("FD_TEST_HF_TOKENIZER_JSON")
	if tokenizerPath == "" {
		tokenizerPath = filepath.Join("..", "..", "tei-models", "deepvk--USER-bge-m3", "tokenizer.json")
	}
	baselinePath := os.Getenv("FD_TEST_TOKENIZER_BASELINE")
	if baselinePath == "" {
		baselinePath = filepath.Join("..", "..", "benchmark-results", "fd-tokenizer-baseline-m012-s01.txt")
	}

	probes := loadTokenBaselineProbes(t, baselinePath)
	rawTexts := map[string]string{
		"contract_question_jurisdiction": "Как определяется подсудность спора по договору поставки между российскими юридическими лицами?",
		"contract_clause_delivery":       "Поставщик обязан передать товар в срок, согласованный сторонами, а покупатель обязан принять товар и оплатить его по цене договора.",
		"labor_question_termination":     "Какие гарантии предоставляются работнику при расторжении трудового договора по инициативе работодателя?",
		"civil_clause_damages":           "Лицо, право которого нарушено, может требовать полного возмещения причинённых ему убытков, если законом или договором не предусмотрено иное.",
		"neutral_russian_reference":      "Москва — крупный научный, культурный и транспортный центр России с развитой городской инфраструктурой.",
	}

	tokenizer, err := newNativeHFTokenizer(tokenizerPath)
	require.NoError(t, err)
	defer tokenizer.Close()

	for _, probe := range probes {
		text, ok := rawTexts[probe.Label]
		require.Truef(t, ok, "missing raw probe mapping for %s", probe.Label)
		encoding, err := tokenizer.Encode(text)
		require.NoError(t, err)
		require.Equalf(t, probe.InputIDs, encoding.InputIDs, "input ids mismatch for %s", probe.Label)
		require.Equalf(t, probe.AttentionMask, encoding.AttentionMask, "attention mask mismatch for %s", probe.Label)
	}
}

func loadTokenBaselineProbes(t *testing.T, path string) []tokenBaselineProbe {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	const marker = "## Token Evidence"
	parts := splitOnce(string(data), marker)
	require.Len(t, parts, 2)
	jsonParts := splitOnce(parts[1], "```json")
	require.Len(t, jsonParts, 2)
	jsonBlockParts := splitOnce(jsonParts[1], "```")
	require.Len(t, jsonBlockParts, 2)
	var probes []tokenBaselineProbe
	require.NoError(t, json.Unmarshal([]byte(jsonBlockParts[0]), &probes))
	require.NotEmpty(t, probes)
	return probes
}

func splitOnce(value, sep string) []string {
	idx := -1
	for i := 0; i+len(sep) <= len(value); i++ {
		if value[i:i+len(sep)] == sep {
			idx = i
			break
		}
	}
	if idx < 0 {
		return []string{value}
	}
	return []string{value[:idx], value[idx+len(sep):]}
}
