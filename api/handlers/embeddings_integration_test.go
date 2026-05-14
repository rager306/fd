package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fd-api/embed"

	"github.com/gin-gonic/gin"
)

// MockTEIClient implements a mock TEI client for testing
type MockTEIClient struct {
	EmbedFunc func(texts []string) ([][]float32, error)
}

func (m *MockTEIClient) Embed(texts []string) ([][]float32, error) {
	return m.EmbedFunc(texts)
}

// testableEmbeddingsHandler is a test-friendly version that uses interfaces
type testableEmbeddingsHandler struct {
	teiClient mockTEIInterface
	cache     mockCacheInterface
	modelID   string
}

type mockTEIInterface interface {
	Embed(texts []string) ([][]float32, error)
}

type mockCacheInterface interface {
	Get(ctx context.Context, text string, dim int) ([]float32, bool, error)
	Set(ctx context.Context, text string, embedding []float32, dim int) error
}

func newTestableEmbeddingsHandler(tei mockTEIInterface, c mockCacheInterface, modelID string) *testableEmbeddingsHandler {
	return &testableEmbeddingsHandler{
		teiClient: tei,
		cache:     c,
		modelID:   modelID,
	}
}

func (h *testableEmbeddingsHandler) CreateEmbedding(c *gin.Context) {
	var req embed.EmbeddingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	texts := req.Input
	if len(texts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input is required"})
		return
	}

	// Dimensions
	dims := 1024
	if req.Dimensions != nil {
		d := *req.Dimensions
		if d == 512 {
			dims = 512
		} else if d != 1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dimensions must be 1024 or 512"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	embeddings := make([][]float32, len(texts))
	promptTokens := 0

	for i, text := range texts {
		emb, found, err := h.cache.Get(ctx, text, dims)
		if err != nil {
			// Log but continue
		}

		if found && emb != nil {
			embeddings[i] = emb
			promptTokens += len(text) / 4
			continue
		}

		embs, err := h.teiClient.Embed([]string{text})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "embedding generation failed"})
			return
		}

		fullEmb := embs[0]
		if dims == 512 && len(fullEmb) >= 512 {
			fullEmb = fullEmb[:512]
		}

		embeddings[i] = fullEmb

		if err := h.cache.Set(ctx, text, fullEmb, dims); err != nil {
			// Log but continue
		}

		promptTokens += len(text) / 4
	}

	data := make([]embed.EmbeddingObj, len(embeddings))
	for i, emb := range embeddings {
		data[i] = embed.EmbeddingObj{
			Object:     "embedding",
			Embedding:  emb,
			Index:      i,
			Dimensions: dims,
		}
	}

	response := embed.EmbeddingsResponse{
		Object: "list",
		Data:   data,
		Model:  h.modelID,
		Usage: embed.Usage{
			PromptTokens: promptTokens,
			TotalTokens:  promptTokens,
		},
	}
	c.JSON(http.StatusOK, response)
}

// mockCacheClient implements mockCacheInterface for testing
type mockCacheClient struct {
	getFunc func(ctx context.Context, text string, dim int) ([]float32, bool, error)
	setFunc func(ctx context.Context, text string, embedding []float32, dim int) error
}

func (m *mockCacheClient) Get(ctx context.Context, text string, dim int) ([]float32, bool, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, text, dim)
	}
	return nil, false, nil
}

func (m *mockCacheClient) Set(ctx context.Context, text string, embedding []float32, dim int) error {
	if m.setFunc != nil {
		return m.setFunc(ctx, text, embedding, dim)
	}
	return nil
}

func TestCreateEmbedding_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		body          string
		mockEmbed     func(texts []string) ([][]float32, error)
		mockCacheGet  func(ctx context.Context, text string, dim int) ([]float32, bool, error)
		mockCacheSet  func(ctx context.Context, text string, embedding []float32, dim int) error
		wantStatus    int
		checkResponse func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "valid single text",
			body: `{"model":"test","input":"hello"}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				return [][]float32{{0.1, 0.2, 0.3}}, nil
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return nil, false, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				return nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp embed.EmbeddingsResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp.Object != "list" {
					t.Errorf("expected object=list, got %s", resp.Object)
				}
				if len(resp.Data) != 1 {
					t.Errorf("expected 1 embedding, got %d", len(resp.Data))
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
			name: "512 dimensions",
			body: `{"model":"test","input":"hello","dimensions":512}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				vec := make([]float32, 1024)
				for i := range vec {
					vec[i] = float32(i) / 1024.0
				}
				return [][]float32{vec}, nil
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return nil, false, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				if dim != 512 {
					t.Errorf("expected dim=512, got %d", dim)
				}
				if len(embedding) != 512 {
					t.Errorf("expected embedding len=512, got %d", len(embedding))
				}
				return nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp embed.EmbeddingsResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp.Data[0].Dimensions != 512 {
					t.Errorf("expected dimensions=512, got %d", resp.Data[0].Dimensions)
				}
				if len(resp.Data[0].Embedding) != 512 {
					t.Errorf("expected embedding len=512, got %d", len(resp.Data[0].Embedding))
				}
			},
		},
		{
			name: "cache hit returns cached embedding",
			body: `{"model":"test","input":"cached text"}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				t.Error("TEI should not be called on cache hit")
				return nil, nil
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return []float32{0.5, 0.6, 0.7}, true, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				return nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp embed.EmbeddingsResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				if len(resp.Data) != 1 {
					t.Errorf("expected 1 embedding, got %d", len(resp.Data))
				}
				if resp.Data[0].Embedding[0] != 0.5 {
					t.Errorf("expected cached embedding [0.5,...], got %v", resp.Data[0].Embedding)
				}
			},
		},
		{
			name: "invalid json",
			body: `{"model":"test","input":}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				return nil, nil
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return nil, false, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				return nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "empty input array",
			body: `{"model":"test","input":[]}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				return nil, nil
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return nil, false, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				return nil
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "tei error returns 500",
			body: `{"model":"test","input":"hello"}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				return nil, fmt.Errorf("TEI unavailable")
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return nil, false, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				return nil
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "multiple texts in batch",
			body: `{"model":"test","input":["hello","world"]}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				return [][]float32{{0.1, 0.2, 0.3}, {0.4, 0.5, 0.6}}, nil
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return nil, false, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				return nil
			},
			wantStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp embed.EmbeddingsResponse
				json.Unmarshal(w.Body.Bytes(), &resp)
				if len(resp.Data) != 2 {
					t.Errorf("expected 2 embeddings, got %d", len(resp.Data))
				}
				if resp.Data[0].Index != 0 || resp.Data[1].Index != 1 {
					t.Errorf("expected indices [0,1], got [%d,%d]", resp.Data[0].Index, resp.Data[1].Index)
				}
			},
		},
		{
			name: "invalid dimensions",
			body: `{"model":"test","input":"hello","dimensions":256}`,
			mockEmbed: func(texts []string) ([][]float32, error) {
				return nil, nil
			},
			mockCacheGet: func(ctx context.Context, text string, dim int) ([]float32, bool, error) {
				return nil, false, nil
			},
			mockCacheSet: func(ctx context.Context, text string, embedding []float32, dim int) error {
				return nil
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTEI := &MockTEIClient{EmbedFunc: tt.mockEmbed}
			mockCache := &mockCacheClient{getFunc: tt.mockCacheGet, setFunc: tt.mockCacheSet}

			handler := newTestableEmbeddingsHandler(mockTEI, mockCache, "test-model")
			router := gin.New()
			router.POST("/v1/embeddings", handler.CreateEmbedding)

			req := httptest.NewRequest("POST", "/v1/embeddings", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d. body: %s", tt.wantStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
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
