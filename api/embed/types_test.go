package embed

import (
	"encoding/json"
	"testing"
)

func TestEmbeddingsRequest_Unmarshal(t *testing.T) {
	// Single string input
	json1 := `{"model":"test","input":"hello"}`
	var req1 EmbeddingsRequest
	if err := json.Unmarshal([]byte(json1), &req1); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(req1.Input) != 1 || req1.Input[0] != "hello" {
		t.Errorf("unexpected: %v", req1.Input)
	}

	// Array input plus OpenAI-compatible optional metadata.
	json2 := `{"model":"test","input":["a","b","c"],"user":"caller-123","priority":"high","encoding_format":"base64"}`
	var req2 EmbeddingsRequest
	if err := json.Unmarshal([]byte(json2), &req2); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(req2.Input) != 3 {
		t.Errorf("expected 3, got %d", len(req2.Input))
	}
	if req2.User == nil || *req2.User != "caller-123" {
		t.Errorf("unexpected user: %v", req2.User)
	}
	if req2.Priority == nil || *req2.Priority != "high" {
		t.Errorf("unexpected priority: %v", req2.Priority)
	}
	if req2.EncodingFormat == nil || *req2.EncodingFormat != EncodingFormatBase64 {
		t.Errorf("unexpected encoding_format: %v", req2.EncodingFormat)
	}
}

func TestEmbeddingsResponse_Marshal(t *testing.T) {
	resp := EmbeddingsResponse{
		Object: "list",
		Data: []EmbeddingObj{
			{Object: "embedding", Embedding: []float32{0.1, -0.2}, Index: 0},
		},
		Model: "test",
		Usage: Usage{PromptTokens: 10, TotalTokens: 10},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var back EmbeddingsResponse
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if back.Object != "list" {
		t.Errorf("expected object=list, got %s", back.Object)
	}
	if len(back.Data) != 1 {
		t.Errorf("expected 1 embedding, got %d", len(back.Data))
	}
	if back.Data[0].Index != 0 {
		t.Errorf("expected index=0, got %d", back.Data[0].Index)
	}
}
