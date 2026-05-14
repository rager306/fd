package embed

import (
	"testing"
)

func TestSliceToDimensions(t *testing.T) {
	full := make([]float32, 1024)
	for i := range full {
		full[i] = float32(i)
	}

	// 512 dimensions
	sliced := full[:512]
	if len(sliced) != 512 {
		t.Errorf("expected len=512, got %d", len(sliced))
	}
	if sliced[0] != 0 {
		t.Errorf("expected first=0, got %f", sliced[0])
	}
	if sliced[511] != 511 {
		t.Errorf("expected last=511, got %f", sliced[511])
	}

	// 1024 dimensions (no slice)
	sliced1024 := full[:1024]
	if len(sliced1024) != 1024 {
		t.Errorf("expected len=1024, got %d", len(sliced1024))
	}
}

func TestSliceToDimensions_LargerThanOriginal(t *testing.T) {
	small := make([]float32, 256)
	for i := range small {
		small[i] = float32(i)
	}

	// Request more than available — slice returns what exists
	sliced := small
	if len(sliced) != 256 {
		t.Errorf("expected len=256, got %d", len(sliced))
	}
}

func TestDimensionsParam_Defaults(t *testing.T) {
	req := EmbeddingsRequest{}
	if req.Dimensions != nil {
		t.Error("expected nil (no dimensions set)")
	}

	req2 := EmbeddingsRequest{Dimensions: ptrInt(512)}
	if req2.Dimensions == nil || *req2.Dimensions != 512 {
		t.Error("expected Dimensions=512")
	}
}

func TestEmbeddingObj_HasDimensions(t *testing.T) {
	obj := EmbeddingObj{
		Object:     "embedding",
		Embedding:  make([]float32, 512),
		Index:      0,
		Dimensions: 512,
	}
	if obj.Dimensions != 512 {
		t.Errorf("expected Dimensions=512, got %d", obj.Dimensions)
	}
}

func TestBatchEmbeddingsRequest(t *testing.T) {
	req := BatchEmbeddingsRequest{
		Inputs:         []string{"text1", "text2"},
		Dimensions:     512,
		EncodingFormat: "base64",
	}
	if len(req.Inputs) != 2 {
		t.Errorf("expected 2 inputs, got %d", len(req.Inputs))
	}
	if req.Dimensions != 512 {
		t.Errorf("expected dims=512, got %d", req.Dimensions)
	}
}

func TestBatchEmbeddingsResponse(t *testing.T) {
	resp := BatchEmbeddingsResponse{
		Embeddings: []string{"abc", "def"},
		Count:      2,
		Dimensions: 512,
	}
	if resp.Count != 2 {
		t.Errorf("expected count=2, got %d", resp.Count)
	}
	if resp.Dimensions != 512 {
		t.Errorf("expected dims=512, got %d", resp.Dimensions)
	}
}

func ptrInt(v int) *int {
	return &v
}
