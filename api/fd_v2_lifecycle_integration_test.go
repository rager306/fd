package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"fd-api/embed"
	"fd-api/handlers"
	"fd-api/lifecycle"
	"fd-api/middleware"

	"github.com/gin-gonic/gin"
)

const lifecycleTestModel = "deepvk/USER-bge-m3"

type lifecycleTestEmbedder struct {
	embedFunc func(context.Context, []string) ([][]float32, error)
}

func (e *lifecycleTestEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return e.embedFunc(ctx, texts)
}

type lifecycleTestCache struct {
	mu    sync.Mutex
	store map[string][]float32
}

func newLifecycleTestCache() *lifecycleTestCache {
	return &lifecycleTestCache{store: make(map[string][]float32)}
}

func (c *lifecycleTestCache) GetIfPresent(_ context.Context, key string, _ int) ([]float32, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	emb, ok := c.store[key]
	return emb, ok
}

func (c *lifecycleTestCache) Set(_ context.Context, key string, _ int, emb []float32) {
	c.mu.Lock()
	defer c.mu.Unlock()
	copyEmb := append([]float32(nil), emb...)
	c.store[key] = copyEmb
}

func (c *lifecycleTestCache) GetManyIfPresent(ctx context.Context, keys []string, dim int) map[int][]float32 {
	hits := make(map[int][]float32, len(keys))
	for i, key := range keys {
		if emb, ok := c.GetIfPresent(ctx, key, dim); ok {
			hits[i] = emb
		}
	}
	return hits
}

func (c *lifecycleTestCache) GetOrLoad(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error) {
	if emb, ok := c.GetIfPresent(ctx, key, dim); ok {
		return emb, nil
	}
	emb, err := loader(ctx)
	if err != nil {
		return nil, err
	}
	c.Set(ctx, key, dim, emb)
	return emb, nil
}

func TestFdV2LifecycleStartupSequence(t *testing.T) {
	state := lifecycle.NewState()
	r := newLifecycleTestRouter(state, 0, staticLifecycleEmbedder())

	assertStatus(t, r, http.MethodGet, "/live", nil, http.StatusOK)
	readyBefore := assertStatus(t, r, http.MethodGet, "/ready", nil, http.StatusServiceUnavailable)
	assertErrorEnvelope(t, readyBefore, handlers.CodeModelNotLoaded, "5")

	state.MarkWarmupDone()
	assertStatus(t, r, http.MethodGet, "/ready", nil, http.StatusOK)
	health := assertStatus(t, r, http.MethodGet, "/health", nil, http.StatusOK)
	var body map[string]any
	if err := json.Unmarshal(health.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal health response: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("health status = %#v, want ok", body["status"])
	}
}

func TestFdV2LifecycleF1ModelNotLoadedThenReady(t *testing.T) {
	state := lifecycle.NewState()
	r := newLifecycleTestRouter(state, 0, staticLifecycleEmbedder())

	beforeWarmup := postLifecycleEmbedding(t, r)
	if beforeWarmup.Code != http.StatusServiceUnavailable {
		t.Fatalf("status before warmup = %d, want %d; body=%s", beforeWarmup.Code, http.StatusServiceUnavailable, beforeWarmup.Body.String())
	}
	assertErrorEnvelope(t, beforeWarmup, handlers.CodeModelNotLoaded, "5")

	state.MarkWarmupDone()
	afterWarmup := postLifecycleEmbedding(t, r)
	if afterWarmup.Code != http.StatusOK {
		t.Fatalf("status after warmup = %d, want %d; body=%s", afterWarmup.Code, http.StatusOK, afterWarmup.Body.String())
	}
}

func TestFdV2LifecycleF2ModelOverloadedThenRecovers(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	started := make(chan struct{})
	release := make(chan struct{})
	var calls atomic.Int64
	embedder := &lifecycleTestEmbedder{embedFunc: func(_ context.Context, _ []string) ([][]float32, error) {
		if calls.Add(1) == 1 {
			close(started)
			<-release
		}
		return staticLifecycleEmbedding(), nil
	}}
	r := newLifecycleTestRouter(state, 1, embedder)

	firstDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		firstDone <- postLifecycleEmbedding(t, r)
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("first embedding request did not start")
	}

	overloaded := postLifecycleEmbedding(t, r)
	if overloaded.Code != http.StatusServiceUnavailable {
		t.Fatalf("overloaded status = %d, want %d; body=%s", overloaded.Code, http.StatusServiceUnavailable, overloaded.Body.String())
	}
	assertErrorEnvelope(t, overloaded, handlers.CodeModelOverloaded, "5")

	close(release)
	select {
	case first := <-firstDone:
		if first.Code != http.StatusOK {
			t.Fatalf("first status = %d, want %d; body=%s", first.Code, http.StatusOK, first.Body.String())
		}
	case <-time.After(time.Second):
		t.Fatal("first request did not finish")
	}

	recovered := postLifecycleEmbedding(t, r)
	if recovered.Code != http.StatusOK {
		t.Fatalf("recovered status = %d, want %d; body=%s", recovered.Code, http.StatusOK, recovered.Body.String())
	}
}

func TestFdV2LifecycleHealthLastInferenceAfterEmbedding(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	r := newLifecycleTestRouter(state, 0, staticLifecycleEmbedder())

	response := postLifecycleEmbedding(t, r)
	if response.Code != http.StatusOK {
		t.Fatalf("embedding status = %d, want %d; body=%s", response.Code, http.StatusOK, response.Body.String())
	}

	health := assertStatus(t, r, http.MethodGet, "/health", nil, http.StatusOK)
	var body handlers.DeepHealthResponse
	if err := json.Unmarshal(health.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal health response: %v", err)
	}
	if body.LastInferenceAt == nil || *body.LastInferenceAt == "" {
		t.Fatalf("last_inference_at missing after successful embedding: %#v", body.LastInferenceAt)
	}
}

func TestFdV2LifecycleF5ShutdownRejectsNewAndDrainsInflight(t *testing.T) {
	state := lifecycle.NewState()
	state.MarkWarmupDone()
	started := make(chan struct{})
	release := make(chan struct{})
	var calls atomic.Int64
	embedder := &lifecycleTestEmbedder{embedFunc: func(_ context.Context, _ []string) ([][]float32, error) {
		if calls.Add(1) == 1 {
			close(started)
			<-release
		}
		return staticLifecycleEmbedding(), nil
	}}
	r := newLifecycleTestRouter(state, 0, embedder)

	firstDone := make(chan *httptest.ResponseRecorder, 1)
	go func() {
		firstDone <- postLifecycleEmbedding(t, r)
	}()
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("in-flight embedding request did not start")
	}

	state.BeginShutdown()
	shuttingDown := postLifecycleEmbedding(t, r)
	if shuttingDown.Code != http.StatusServiceUnavailable {
		t.Fatalf("shutdown status = %d, want %d; body=%s", shuttingDown.Code, http.StatusServiceUnavailable, shuttingDown.Body.String())
	}
	assertErrorEnvelope(t, shuttingDown, handlers.CodeShuttingDown, "30")

	close(release)
	select {
	case first := <-firstDone:
		if first.Code != http.StatusOK {
			t.Fatalf("in-flight status = %d, want %d; body=%s", first.Code, http.StatusOK, first.Body.String())
		}
	case <-time.After(time.Second):
		t.Fatal("in-flight request did not drain")
	}
}

func newLifecycleTestRouter(state *lifecycle.State, maxInFlight int64, embedder embed.Embedder) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	embedHandler := handlers.NewEmbeddingsHandler(embedder, newLifecycleTestCache(), lifecycleTestModel, logger)
	r.GET("/live", handlers.NewLiveHandler())
	r.GET("/ready", handlers.NewReadyHandler(state))
	r.GET("/health", handlers.NewHealthHandlerWithState(&handlers.RuntimeHealth{Backend: "test", Model: lifecycleTestModel}, state))
	r.POST("/v1/embeddings", middleware.ValidateEmbeddingsRequest(), middleware.LifecycleGateWithCapacity(state, maxInFlight), embedHandler.CreateEmbedding)
	return r
}

func postLifecycleEmbedding(t *testing.T, r http.Handler) *httptest.ResponseRecorder {
	t.Helper()
	return assertStatus(t, r, http.MethodPost, "/v1/embeddings", []byte(`{"model":"test","input":"hello"}`), 0)
}

func assertStatus(t *testing.T, r http.Handler, method, path string, body []byte, wantStatus int) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if wantStatus != 0 && w.Code != wantStatus {
		t.Fatalf("%s %s status = %d, want %d; body=%s", method, path, w.Code, wantStatus, w.Body.String())
	}
	return w
}

func assertErrorEnvelope(t *testing.T, w *httptest.ResponseRecorder, code, retryAfter string) {
	t.Helper()
	if got := w.Header().Get("Retry-After"); got != retryAfter {
		t.Fatalf("Retry-After = %q, want %q", got, retryAfter)
	}
	var body handlers.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	if body.Error.Code != code {
		t.Fatalf("error.code = %q, want %q", body.Error.Code, code)
	}
	if body.Error.Type != handlers.TypeOverloadedError {
		t.Fatalf("error.type = %q, want %q", body.Error.Type, handlers.TypeOverloadedError)
	}
}

func staticLifecycleEmbedder() *lifecycleTestEmbedder {
	return &lifecycleTestEmbedder{embedFunc: func(_ context.Context, _ []string) ([][]float32, error) {
		return staticLifecycleEmbedding(), nil
	}}
}

func staticLifecycleEmbedding() [][]float32 {
	emb := make([]float32, 1024)
	emb[0] = 1
	return [][]float32{emb}
}
