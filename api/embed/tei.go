package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type TEIClient struct {
	baseURL    string
	modelID    string
	httpClient *http.Client
}

func NewTEIClient(baseURL, modelID string, client *http.Client) *TEIClient {
	return &TEIClient{
		baseURL:    baseURL,
		modelID:    modelID,
		httpClient: client,
	}
}

// TEI OpenAI-compatible request
type teiRequest struct {
	Input    string `json:"input"`
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

func (c *TEIClient) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	// Single text — use first element
	text := texts[0]
	reqBody := teiRequest{
		Input:    text,
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
