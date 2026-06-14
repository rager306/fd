package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestV1BatchHandlerCreatesEmbeddingsForBatches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	embedder := &mockEmbedder{embedFunc: func(_ context.Context, texts []string) ([][]float32, error) {
		vectors := make([][]float32, len(texts))
		for i := range texts {
			vectors[i] = []float32{float32(len(texts[i])), float32(i)}
		}
		return vectors, nil
	}}
	router := gin.New()
	handler := NewV1BatchHandler(embedder, &mockEmbeddingCache{}, testLogger())
	router.POST("/v1/batch", handler.CreateBatch)

	w := postJSON(router, "/v1/batch", `{"batches":[["a","b","c","d"],["e","f","g","h"]]}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	var resp v1BatchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(resp.Batches) != 2 {
		t.Fatalf("outer batches = %d, want 2", len(resp.Batches))
	}
	for i, batch := range resp.Batches {
		if len(batch) != 4 {
			t.Fatalf("batch %d embeddings = %d, want 4", i, len(batch))
		}
		for j, vector := range batch {
			if len(vector) != 2 {
				t.Fatalf("batch %d embedding %d length = %d, want 2", i, j, len(vector))
			}
		}
	}
}

func TestV1BatchHandlerRejectsOversizedInnerBatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewV1BatchHandler(&mockEmbedder{}, &mockEmbeddingCache{}, testLogger())
	router.POST("/v1/batch", handler.CreateBatch)

	inputs := make([]string, maxV1BatchInputs+1)
	for i := range inputs {
		inputs[i] = "x"
	}
	body, err := json.Marshal(gin.H{"batches": [][]string{inputs}})
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	w := postJSON(router, "/v1/batch", string(body))
	assertV1BatchError(t, w, CodeBatchTooLarge, http.StatusRequestEntityTooLarge, "batches[0]")
}

func TestV1BatchHandlerRejectsEmptyBatches(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewV1BatchHandler(&mockEmbedder{}, &mockEmbeddingCache{}, testLogger())
	router.POST("/v1/batch", handler.CreateBatch)

	w := postJSON(router, "/v1/batch", `{"batches":[]}`)
	assertV1BatchError(t, w, CodeInputRequired, http.StatusBadRequest, "batches")
}

func TestV1BatchHandlerRejectsTooLongInputBeforeEmbedder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	embedder := &mockEmbedder{embedFunc: func(context.Context, []string) ([][]float32, error) {
		t.Fatal("embedder should not be called for invalid input")
		return nil, nil
	}}
	router := gin.New()
	handler := NewV1BatchHandler(embedder, &mockEmbeddingCache{}, testLogger())
	router.POST("/v1/batch", handler.CreateBatch)

	body, err := json.Marshal(gin.H{"batches": [][]string{{string(make([]byte, maxBatchInputChars+1))}}})
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	w := postJSON(router, "/v1/batch", string(body))
	assertV1BatchError(t, w, CodeInputTooLong, http.StatusRequestEntityTooLarge, "batches[0][0]")
	if embedder.calls != 0 {
		t.Fatalf("embedder calls = %d, want 0", embedder.calls)
	}
}

func assertV1BatchError(t *testing.T, w *httptest.ResponseRecorder, wantCode string, wantStatus int, wantParam string) {
	t.Helper()
	if w.Code != wantStatus {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, wantStatus, w.Body.String())
	}
	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal error response: %v; body=%s", err, w.Body.String())
	}
	if resp.Error.Code != wantCode {
		t.Fatalf("error.code = %q, want %q", resp.Error.Code, wantCode)
	}
	if resp.Error.Param != wantParam {
		t.Fatalf("error.param = %q, want %q", resp.Error.Param, wantParam)
	}
}
