package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"fd-api/embed"

	"github.com/gin-gonic/gin"
)

type mockEmbedder struct {
	embedFunc func(ctx context.Context, texts []string) ([][]float32, error)
	calls     int
}

func (m *mockEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	m.calls++
	return m.embedFunc(ctx, texts)
}

type mockEmbeddingCache struct {
	getOrLoadFunc func(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error)
}

func (m *mockEmbeddingCache) GetOrLoad(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error) {
	if m.getOrLoadFunc != nil {
		return m.getOrLoadFunc(ctx, key, dim, loader)
	}
	return loader(ctx)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func bufferLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

func postJSON(router *gin.Engine, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestCreateEmbedding_ProductionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		body          string
		embedFunc     func(ctx context.Context, texts []string) ([][]float32, error)
		cacheFunc     func(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error)
		wantStatus    int
		checkResponse func(*testing.T, *httptest.ResponseRecorder, *mockEmbedder)
	}{
		{
			name: "valid single text",
			body: `{"model":"test","input":"hello"}`,
			embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
				return [][]float32{{0.1, 0.2, 0.3}}, nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, _ *mockEmbedder) {
				var resp embed.EmbeddingsResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				if resp.Object != "list" {
					t.Errorf("expected object=list, got %s", resp.Object)
				}
				if len(resp.Data) != 1 {
					t.Fatalf("expected 1 embedding, got %d", len(resp.Data))
				}
				if resp.Data[0].Index != 0 {
					t.Errorf("expected index=0, got %d", resp.Data[0].Index)
				}
				if resp.Data[0].Dimensions != 1024 {
					t.Errorf("expected dimensions=1024, got %d", resp.Data[0].Dimensions)
				}
			},
		},
		{
			name: "512 dimensions slices embedding",
			body: `{"model":"test","input":"hello","dimensions":512}`,
			embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
				vec := make([]float32, 1024)
				for i := range vec {
					vec[i] = float32(i) / 1024.0
				}
				return [][]float32{vec}, nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, _ *mockEmbedder) {
				var resp embed.EmbeddingsResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				if resp.Data[0].Dimensions != 512 {
					t.Errorf("expected dimensions=512, got %d", resp.Data[0].Dimensions)
				}
				if len(resp.Data[0].Embedding) != 512 {
					t.Errorf("expected embedding len=512, got %d", len(resp.Data[0].Embedding))
				}
			},
		},
		{
			name: "cache hit does not call embedder",
			body: `{"model":"test","input":"cached text"}`,
			embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
				t.Fatal("embedder should not be called on cache hit")
				return nil, nil
			},
			cacheFunc: func(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error) {
				return []float32{0.5, 0.6, 0.7}, nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, embedder *mockEmbedder) {
				if embedder.calls != 0 {
					t.Fatalf("embedder calls=%d, want 0", embedder.calls)
				}
				var resp embed.EmbeddingsResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				if resp.Data[0].Embedding[0] != 0.5 {
					t.Errorf("expected cached embedding [0.5,...], got %v", resp.Data[0].Embedding)
				}
			},
		},
		{
			name:       "invalid json",
			body:       `{"model":"test","input":}`,
			embedFunc:  func(ctx context.Context, texts []string) ([][]float32, error) { return nil, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty input array",
			body:       `{"model":"test","input":[]}`,
			embedFunc:  func(ctx context.Context, texts []string) ([][]float32, error) { return nil, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "tei error returns 500",
			body: `{"model":"test","input":"hello"}`,
			embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
				return nil, errors.New("TEI unavailable")
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid dimensions",
			body:       `{"model":"test","input":"hello","dimensions":256}`,
			embedFunc:  func(ctx context.Context, texts []string) ([][]float32, error) { return nil, nil },
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embedder := &mockEmbedder{embedFunc: tt.embedFunc}
			cache := &mockEmbeddingCache{getOrLoadFunc: tt.cacheFunc}
			handler := NewEmbeddingsHandler(embedder, cache, "test-model", testLogger())
			router := gin.New()
			router.POST("/v1/embeddings", handler.CreateEmbedding)

			w := postJSON(router, "/v1/embeddings", tt.body)
			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d. body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
			if tt.checkResponse != nil {
				tt.checkResponse(t, w, embedder)
			}
		})
	}
}

func TestCreateBatchEmbeddings_Validation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	embedder := &mockEmbedder{embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
		vec := make([]float32, 1024)
		return [][]float32{vec}, nil
	}}
	cache := &mockEmbeddingCache{}
	handler := NewBatchHandler(embedder, cache, "test-model", testLogger())
	router := gin.New()
	router.POST("/embeddings/batch", handler.CreateBatchEmbeddings)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid default request",
			body:       `{"inputs":["hello"]}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "valid 512d float request",
			body:       `{"inputs":["hello"],"dimensions":512,"encoding_format":"float"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid dimensions rejected",
			body:       `{"inputs":["hello"],"dimensions":256}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid encoding format rejected",
			body:       `{"inputs":["hello"],"encoding_format":"hex"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty inputs rejected",
			body:       `{"inputs":[]}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := postJSON(router, "/embeddings/batch", tt.body)
			if w.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d. body: %s", tt.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestCreateBatchEmbeddings_Base64Response(t *testing.T) {
	gin.SetMode(gin.TestMode)
	embedder := &mockEmbedder{embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
		return [][]float32{{1, 2, 3, 4}}, nil
	}}
	handler := NewBatchHandler(embedder, &mockEmbeddingCache{}, "test-model", testLogger())
	router := gin.New()
	router.POST("/embeddings/batch", handler.CreateBatchEmbeddings)

	w := postJSON(router, "/embeddings/batch", `{"inputs":["hello"]}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}

	var resp embed.BatchEmbeddingsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Count != 1 || resp.Dimensions != 1024 {
		t.Fatalf("unexpected response metadata: %+v", resp)
	}
	if _, err := base64.StdEncoding.DecodeString(resp.Embeddings[0]); err != nil {
		t.Fatalf("embedding is not valid base64: %v", err)
	}
}

func TestCreateEmbedding_SuccessDoesNotEmitInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logs bytes.Buffer
	embedder := &mockEmbedder{embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
		return [][]float32{{0.1, 0.2, 0.3}}, nil
	}}
	handler := NewEmbeddingsHandler(embedder, &mockEmbeddingCache{}, "test-model", bufferLogger(&logs))
	router := gin.New()
	router.POST("/v1/embeddings", handler.CreateEmbedding)

	w := postJSON(router, "/v1/embeddings", `{"model":"test","input":"hello"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
	if logs.Len() != 0 {
		t.Fatalf("successful embedding request emitted info logs: %s", logs.String())
	}
}

func TestCreateBatchEmbeddings_SuccessDoesNotEmitInfo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var logs bytes.Buffer
	embedder := &mockEmbedder{embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
		return [][]float32{{1, 2, 3, 4}}, nil
	}}
	handler := NewBatchHandler(embedder, &mockEmbeddingCache{}, "test-model", bufferLogger(&logs))
	router := gin.New()
	router.POST("/embeddings/batch", handler.CreateBatchEmbeddings)

	w := postJSON(router, "/embeddings/batch", `{"inputs":["hello"]}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
	}
	if logs.Len() != 0 {
		t.Fatalf("successful batch request emitted info logs: %s", logs.String())
	}
}

func TestEmbeddingsRequest_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single string input",
			input:    `{"model":"test","input":"hello"}`,
			expected: []string{"hello"},
		},
		{
			name:     "array string input",
			input:    `{"model":"test","input":["hello","world"]}`,
			expected: []string{"hello", "world"},
		},
		{
			name:     "dimensions field parsed",
			input:    `{"model":"test","input":"hello","dimensions":512}`,
			expected: []string{"hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req embed.EmbeddingsRequest
			err := json.Unmarshal([]byte(tt.input), &req)
			if err != nil {
				t.Fatalf("unmarshal error: %v", err)
			}
			if len(req.Input) != len(tt.expected) {
				t.Errorf("expected %d inputs, got %d", len(tt.expected), len(req.Input))
			}
			for i, exp := range tt.expected {
				if req.Input[i] != exp {
					t.Errorf("expected input[%d]=%q, got %q", i, exp, req.Input[i])
				}
			}
		})
	}
}
