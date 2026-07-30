package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"fd-api/embed"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	helloInputJSON = `{"model":"test","input":"hello"}`
	helloText      = "hello"
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
	getOrLoadFunc        func(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error)
	getManyIfPresentFunc func(ctx context.Context, keys []string, dim int) map[int][]float32
	getIfPresentCalls    int
	getManyCalls         int
	// In-memory store for GetIfPresent / Set pair (mirrors the
	// production TieredCache shape closely enough for handler tests).
	store map[string][]float32
}

func (m *mockEmbeddingCache) GetIfPresent(ctx context.Context, key string, dim int) ([]float32, bool) {
	m.getIfPresentCalls++
	if m.store == nil {
		return nil, false
	}
	v, ok := m.store[fmt.Sprintf("%s:d%d", key, dim)]
	return v, ok
}

func (m *mockEmbeddingCache) GetManyIfPresent(ctx context.Context, keys []string, dim int) map[int][]float32 {
	m.getManyCalls++
	if m.getManyIfPresentFunc != nil {
		return m.getManyIfPresentFunc(ctx, keys, dim)
	}
	hits := make(map[int][]float32)
	for i, key := range keys {
		if v, ok := m.GetIfPresent(ctx, key, dim); ok {
			hits[i] = v
		}
	}
	return hits
}

func (m *mockEmbeddingCache) Set(ctx context.Context, key string, dim int, emb []float32) {
	if m.store == nil {
		m.store = map[string][]float32{}
	}
	if len(emb) < dim {
		return
	}
	m.store[fmt.Sprintf("%s:d%d", key, dim)] = emb[:dim]
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

func postJSON(router *gin.Engine, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestCreateEmbeddingUsesBatchedCachePeek(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := &mockEmbeddingCache{
		getManyIfPresentFunc: func(ctx context.Context, keys []string, dim int) map[int][]float32 {
			if dim != 1024 {
				t.Fatalf("dim = %d, want 1024", dim)
			}
			if fmt.Sprint(keys) != "[a b c]" {
				t.Fatalf("keys = %v, want [a b c]", keys)
			}
			vec := make([]float32, 1024)
			vec[0] = 20
			return map[int][]float32{1: vec}
		},
	}
	var seen [][]string
	embedder := &mockEmbedder{embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
		seen = append(seen, append([]string(nil), texts...))
		vectors := make([][]float32, len(texts))
		for i, text := range texts {
			vec := make([]float32, 1024)
			vec[0] = float32(len(text))
			vec[1] = float32(i)
			vectors[i] = vec
		}
		return vectors, nil
	}}
	handler := NewEmbeddingsHandler(embedder, cache, "test-model", testLogger())
	router := gin.New()
	router.POST("/v1/embeddings", handler.CreateEmbedding)

	w := postJSON(router, "/v1/embeddings", `{"model":"test","input":["a","b","c"]}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if cache.getManyCalls != 1 {
		t.Fatalf("GetManyIfPresent calls = %d, want 1", cache.getManyCalls)
	}
	if cache.getIfPresentCalls != 0 {
		t.Fatalf("GetIfPresent calls = %d, want 0", cache.getIfPresentCalls)
	}
	if embedder.calls != 1 {
		t.Fatalf("embedder calls = %d, want 1", embedder.calls)
	}
	if fmt.Sprint(seen) != "[[a c]]" {
		t.Fatalf("embedder texts = %v, want [[a c]]", seen)
	}
}

//nolint:gocyclo // table-driven production integration coverage intentionally exercises many request/error/cache paths in one matrix.
func TestCreateEmbedding_ProductionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		body          string
		embedFunc     func(ctx context.Context, texts []string) ([][]float32, error)
		cacheFunc     func(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error)
		cacheSetup    func(c *mockEmbeddingCache) // pre-populate cache for GetIfPresent path
		wantStatus    int
		checkResponse func(*testing.T, *httptest.ResponseRecorder, *mockEmbedder)
	}{
		{
			name: "valid single text",
			body: helloInputJSON,
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
			name: "base64 encoding format response",
			body: `{"model":"test","input":"hello","encoding_format":"base64"}`,
			embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
				return [][]float32{{1, 2, 3, 4}}, nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, _ *mockEmbedder) {
				var resp embed.EmbeddingsResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				embedding, ok := resp.Data[0].Embedding.(string)
				if !ok {
					t.Fatalf("expected base64 string embedding, got %T", resp.Data[0].Embedding)
				}
				decoded, err := base64.StdEncoding.DecodeString(embedding)
				if err != nil {
					t.Fatalf("decode base64 embedding: %v", err)
				}
				if len(decoded) != 16 {
					t.Fatalf("decoded length = %d, want 16 bytes", len(decoded))
				}
			},
		},
		{
			name: "priority and user are accepted",
			body: `{"model":"test","input":"hello","priority":"high","user":"caller-123"}`,
			embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
				return [][]float32{{0.1, 0.2, 0.3}}, nil
			},
			wantStatus: http.StatusOK,
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
				// JSON round-trip via interface{} yields []interface{} of
				// float64 (Go's json package has no native float32). Re-marshal
				// to []float32 to verify the wire shape is correct.
				raw, ok := resp.Data[0].Embedding.([]interface{})
				if !ok {
					t.Fatalf("expected []interface{} embedding after JSON round-trip, got %T", resp.Data[0].Embedding)
				}
				if len(raw) != 512 {
					t.Errorf("expected embedding len=512, got %d", len(raw))
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
			// Pre-populate the mock cache with a 1024-dim vector, so the
			// new GetIfPresent-based handler code path short-circuits
			// without hitting the embedder. Mock Set requires len(emb)
			// >= dim (mirrors production TieredCache.Set).
			cacheSetup: func(c *mockEmbeddingCache) {
				cached := make([]float32, 1024)
				cached[0] = 0.5
				cached[1] = 0.6
				cached[2] = 0.7
				c.Set(context.Background(), "cached text", 1024, cached)
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, embedder *mockEmbedder) {
				if got := w.Header().Get(HeaderCache); got != cacheHit {
					t.Fatalf("X-Cache = %q, want %q", got, cacheHit)
				}
				if embedder.calls != 0 {
					t.Fatalf("embedder calls=%d, want 0", embedder.calls)
				}
				var resp embed.EmbeddingsResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				raw, ok := resp.Data[0].Embedding.([]interface{})
				if !ok {
					t.Fatalf("expected []interface{} embedding, got %T", resp.Data[0].Embedding)
				}
				if len(raw) == 0 {
					t.Fatalf("expected non-empty embedding")
				}
				first, ok := raw[0].(float64)
				if !ok {
					t.Fatalf("expected float64 element, got %T", raw[0])
				}
				if first != 0.5 {
					t.Errorf("expected first element=0.5, got %v", first)
				}
				if len(raw) != 1024 {
					t.Errorf("expected embedding len=1024, got %d", len(raw))
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
			name:       "invalid priority",
			body:       `{"model":"test","input":"hello","priority":"urgent"}`,
			embedFunc:  func(ctx context.Context, texts []string) ([][]float32, error) { return nil, nil },
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "tei error returns 500",
			body: helloInputJSON,
			embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
				return nil, errors.New("TEI unavailable")
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "invalid dimensions",
			body:       `{"model":"test","input":helloText,"dimensions":256}`,
			embedFunc:  func(ctx context.Context, texts []string) ([][]float32, error) { return nil, nil },
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			embedder := &mockEmbedder{embedFunc: tt.embedFunc}
			cache := &mockEmbeddingCache{getOrLoadFunc: tt.cacheFunc}
			if tt.cacheSetup != nil {
				tt.cacheSetup(cache)
			}
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

func TestCreateBatchEmbeddingsUsesSingleEmbedCallForMisses(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var seen [][]string
	embedder := &mockEmbedder{embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
		seen = append(seen, append([]string(nil), texts...))
		vectors := make([][]float32, len(texts))
		for i, text := range texts {
			vec := make([]float32, 1024)
			vec[0] = float32(len(text))
			vec[1] = float32(i)
			vectors[i] = vec
		}
		return vectors, nil
	}}
	handler := NewBatchHandler(embedder, &mockEmbeddingCache{}, "test-model", testLogger())
	router := gin.New()
	router.POST("/embeddings/batch", handler.CreateBatchEmbeddings)

	w := postJSON(router, "/embeddings/batch", `{"inputs":["a","b","c","d"]}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if embedder.calls != 1 {
		t.Fatalf("embedder calls = %d, want 1", embedder.calls)
	}
	wantSeen := [][]string{{"a", "b", "c", "d"}}
	if !reflect.DeepEqual(seen, wantSeen) {
		t.Fatalf("embedder texts = %#v, want %#v", seen, wantSeen)
	}

	w = postJSON(router, "/embeddings/batch", `{"inputs":["a","b","c","d"]}`)
	if w.Code != http.StatusOK {
		t.Fatalf("second status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if embedder.calls != 1 {
		t.Fatalf("embedder calls after cache hit = %d, want 1", embedder.calls)
	}
}

func TestCreateBatchEmbeddingsRejectsTooLongInputBeforeEmbedder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	embedder := &mockEmbedder{embedFunc: func(context.Context, []string) ([][]float32, error) {
		t.Fatal("embedder should not be called for invalid input")
		return nil, nil
	}}
	handler := NewBatchHandler(embedder, &mockEmbeddingCache{}, "test-model", testLogger())
	router := gin.New()
	router.POST("/embeddings/batch", handler.CreateBatchEmbeddings)

	body, err := json.Marshal(gin.H{"inputs": []string{string(make([]byte, maxBatchInputChars+1))}})
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	w := postJSON(router, "/embeddings/batch", string(body))
	assertV1BatchError(t, w, CodeInputTooLong, http.StatusRequestEntityTooLarge, "inputs[0]")
	if embedder.calls != 0 {
		t.Fatalf("embedder calls = %d, want 0", embedder.calls)
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
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp), "unmarshal response")
	assert.Equal(t, 1, resp.Count)
	assert.Equal(t, 1024, resp.Dimensions)
	_, err := base64.StdEncoding.DecodeString(resp.Embeddings[0])
	assert.NoError(t, err, "embedding should be valid base64")
}

func TestCreateEmbedding_ModelRequestFieldIsCompatibilityMetadata(t *testing.T) {
	gin.SetMode(gin.TestMode)

	embedder := &mockEmbedder{embedFunc: func(ctx context.Context, texts []string) ([][]float32, error) {
		return [][]float32{{0.1, 0.2, 0.3}}, nil
	}}
	handler := NewEmbeddingsHandler(embedder, &mockEmbeddingCache{}, "configured-model", testLogger())
	router := gin.New()
	router.POST("/v1/embeddings", handler.CreateEmbedding)

	w := postJSON(router, "/v1/embeddings", `{"model":"client-placeholder","input":"hello"}`)
	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())

	var resp embed.EmbeddingsResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp), "unmarshal response")
	assert.Equal(t, "configured-model", resp.Model, "response model should be authoritative")
	assert.Equal(t, 1, embedder.calls, "mismatched request model should not block embedding generation")
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

	w := postJSON(router, "/v1/embeddings", helloInputJSON)
	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	assert.Empty(t, logs.String(), "successful embedding request should not emit info logs")
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
	require.Equal(t, http.StatusOK, w.Code, "body: %s", w.Body.String())
	assert.Empty(t, logs.String(), "successful batch request should not emit info logs")
}

func TestEmbeddingsRequest_Unmarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single string input",
			input:    helloInputJSON,
			expected: []string{helloText},
		},
		{
			name:     "array string input",
			input:    `{"model":"test","input":["hello","world"]}`,
			expected: []string{helloText, "world"},
		},
		{
			name:     "dimensions field parsed",
			input:    `{"model":"test","input":"hello","dimensions":512}`,
			expected: []string{helloText},
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
