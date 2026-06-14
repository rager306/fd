package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// TEIClient calls a Text Embeddings Inference (TEI) OpenAI-compatible
// endpoint and converts its nested response into fd's [][]float32 interface.
type TEIClient struct {
	baseURL    string
	modelID    string
	httpClient *http.Client
}

// NewTEIClient returns a TEIClient for baseURL/modelID using the supplied
// HTTP client. The caller owns the HTTP client's timeout and transport settings.
func NewTEIClient(baseURL, modelID string, client *http.Client) *TEIClient {
	return &TEIClient{
		baseURL:    baseURL,
		modelID:    modelID,
		httpClient: client,
	}
}

// TEI OpenAI-compatible request
type teiRequest struct {
	Input    any    `json:"input"` // string OR []string (OpenAI spec)
	Model    string `json:"model"`
	Truncate bool   `json:"truncate"`
}

// TEI OpenAI-compatible response — NESTED structure (data[0].embedding)
type teiResponse struct {
	Object string            `json:"object"`
	Data   []teiEmbeddingObj `json:"data"`
	Model  string            `json:"model"`
	Usage  Usage             `json:"usage"`
}

type teiEmbeddingObj struct {
	Embedding []float32 `json:"embedding"`
}

// Embed sends texts to TEI's OpenAI-compatible embeddings endpoint and
// returns the embedding vectors in response order.
func (c *TEIClient) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	// Send the FULL slice as TEI input. TEI accepts both single string and
	// []string via the OpenAI-compatible endpoint (the spec's "input" field
	// is defined as either). Previously this code used texts[0] only,
	// which silently dropped batch elements when the caller passed >1
	// input — see S01-R01 for the bug pattern. Fixed 2026-06 when fd v2
	// raised maxBatchSize from 32 to 128, which made the per-item
	// chunking in the handler start sending multi-element chunks.
	reqBody := teiRequest{
		Input:    texts, // marshal as JSON []string when slice, single string when len==1
		Model:    c.modelID,
		Truncate: true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/embeddings", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TEI returned status %d", resp.StatusCode)
	}

	var result teiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf("TEI returned empty data")
	}

	embeddings := make([][]float32, len(result.Data))
	for i, d := range result.Data {
		embeddings[i] = d.Embedding
	}

	return embeddings, nil
}
