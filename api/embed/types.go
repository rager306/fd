package embed

import "encoding/json"

// OpenAI-compatible request
type EmbeddingsRequest struct {
	Model      string   `json:"model"`
	Input      []string `json:"input"` // slice for batch, filled from string or []string
	Dimensions *int     `json:"dimensions"` // pointer: nil=1024 (default), 512, or explicitly 1024
}

// UnmarshalJSON implements custom JSON unmarshaling to handle both string and []string input
func (r *EmbeddingsRequest) UnmarshalJSON(data []byte) error {
	type rawRequest struct {
		Model      string          `json:"model"`
		Input      json.RawMessage `json:"input"`
		Dimensions *int            `json:"dimensions"`
	}
	var raw rawRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Model = raw.Model
	r.Dimensions = raw.Dimensions

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

// OpenAI-compatible response
type EmbeddingsResponse struct {
	Object string         `json:"object"`
	Data   []EmbeddingObj `json:"data"`
	Model  string         `json:"model"`
	Usage  Usage          `json:"usage"`
}

type EmbeddingObj struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
	Dimensions int      `json:"dimensions"` // 1024 or 512
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// Batch request (internal endpoint for FalkorDB)
type BatchEmbeddingsRequest struct {
	Inputs         []string `json:"inputs"`
	Dimensions     int      `json:"dimensions"`      // 1024 or 512, default 1024
	EncodingFormat string   `json:"encoding_format"` // "base64" or "float"
}

// Batch response
type BatchEmbeddingsResponse struct {
	Embeddings []string `json:"embeddings"` // base64 encoded
	Count      int      `json:"count"`
	Dimensions int      `json:"dimensions"`
}
