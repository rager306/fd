package embed

import (
	"context"
	"encoding/json"
)

// Embedder is the minimal inference interface shared by handlers and lifecycle warmup.
type Embedder interface {
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}

// EmbeddingsRequest is the OpenAI-compatible /v1/embeddings request body.
// Input accepts either a single string or []string via custom JSON unmarshaling.
type EmbeddingsRequest struct {
	Model          string   `json:"model"`
	Input          []string `json:"input"`           // slice for batch, filled from string or []string
	Dimensions     *int     `json:"dimensions"`      // pointer: nil=1024 (default), 512, or explicitly 1024
	EncodingFormat *string  `json:"encoding_format"` // pointer: nil=float (default), "float", or "base64"
	User           *string  `json:"user"`            // optional caller identifier for compatibility/observability
	Priority       *string  `json:"priority"`        // optional priority: low, normal, or high
}

// UnmarshalJSON implements custom JSON unmarshaling to handle both string and []string input.
// When the input field is absent (raw.Input is nil), r.Input stays nil so the
// middleware can distinguish "field absent" from "field present but empty"
// and emit input_required vs accepting the empty value.
func (r *EmbeddingsRequest) UnmarshalJSON(data []byte) error {
	type rawRequest struct {
		Model          string          `json:"model"`
		Input          json.RawMessage `json:"input"`
		Dimensions     *int            `json:"dimensions"`
		EncodingFormat *string         `json:"encoding_format"`
		User           *string         `json:"user"`
		Priority       *string         `json:"priority"`
	}
	var raw rawRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Model = raw.Model
	r.Dimensions = raw.Dimensions
	r.EncodingFormat = raw.EncodingFormat
	r.User = raw.User
	r.Priority = raw.Priority

	// Field absent — leave r.Input as nil; validation middleware emits input_required.
	if len(raw.Input) == 0 {
		return nil
	}
	// Try []string first
	if err := json.Unmarshal(raw.Input, &r.Input); err == nil {
		return nil
	}
	// Fall back to single String
	var single string
	if err := json.Unmarshal(raw.Input, &single); err != nil {
		return err
	}
	r.Input = []string{single}
	return nil
}

// EmbeddingsResponse is the OpenAI-compatible /v1/embeddings response body.
type EmbeddingsResponse struct {
	Object string         `json:"object"`
	Data   []EmbeddingObj `json:"data"`
	Model  string         `json:"model"`
	Usage  Usage          `json:"usage"`
}

// EmbeddingObj is one item in the OpenAI-compatible response data array.
type EmbeddingObj struct {
	Object string `json:"object"`
	// Embedding carries the vector. When EncodingFormat=float (default),
	// this is []float32 marshaled as a JSON array. When EncodingFormat=base64,
	// this is a base64-encoded string of the float32 LE bytes.
	// Use EmbeddingObj.SetVector / SetBase64 to populate.
	Embedding  any `json:"embedding"`
	Index      int `json:"index"`
	Dimensions int `json:"dimensions"` // 1024 or 512
}

// SetVector stores a float32 slice. JSON marshals it as a numeric array.
func (e *EmbeddingObj) SetVector(v []float32) { e.Embedding = v }

// SetBase64 stores a base64-encoded float32 LE byte string.
func (e *EmbeddingObj) SetBase64(s string) { e.Embedding = s }

// Usage reports token accounting fields in the OpenAI-compatible response.
type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// BatchEmbeddingsRequest is the legacy internal /embeddings/batch request
// shape used by FalkorDB callers.
type BatchEmbeddingsRequest struct {
	Inputs         []string `json:"inputs"`
	Dimensions     int      `json:"dimensions"`      // 1024 or 512, default 1024
	EncodingFormat string   `json:"encoding_format"` // "base64" or "float"
}

// BatchEmbeddingsResponse is the legacy internal /embeddings/batch response
// shape with base64-encoded embedding vectors.
type BatchEmbeddingsResponse struct {
	Embeddings []string `json:"embeddings"` // base64 encoded
	Count      int      `json:"count"`
	Dimensions int      `json:"dimensions"`
}
