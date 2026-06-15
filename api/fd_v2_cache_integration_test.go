package main

import (
	"context"
	"encoding/binary"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"fd-api/cache"
	"fd-api/handlers"
	"fd-api/middleware"
	"fd-api/observability"

	"github.com/gin-gonic/gin"
)

const testVectorCacheTTL = time.Hour

type localVectorCache struct {
	local   *cache.LocalCache
	metrics *observability.Metrics
}

func newLocalVectorCache(maxSize int, metrics *observability.Metrics) *localVectorCache {
	return &localVectorCache{local: cache.NewLocalCache(maxSize, 0), metrics: metrics}
}

func (c *localVectorCache) GetIfPresent(ctx context.Context, key string, dim int) ([]float32, bool) {
	data, ok := c.local.Get(ctx, key)
	if !ok {
		c.observeCacheResult("miss")
		return nil, false
	}
	vec, ok := decodeTestVector(data, dim)
	if !ok {
		c.observeCacheResult("miss")
		return nil, false
	}
	c.observeCacheResult("hit")
	return vec, true
}

func (c *localVectorCache) GetManyIfPresent(ctx context.Context, keys []string, dim int) map[int][]float32 {
	hits := make(map[int][]float32)
	for i, key := range keys {
		if vec, ok := c.GetIfPresent(ctx, key, dim); ok {
			hits[i] = vec
		}
	}
	return hits
}

func (c *localVectorCache) Set(ctx context.Context, key string, _ int, emb []float32) {
	c.local.Set(ctx, key, encodeTestVector(emb), testVectorCacheTTL)
}

func (c *localVectorCache) GetOrLoad(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error) {
	if vec, ok := c.GetIfPresent(ctx, key, dim); ok {
		return vec, nil
	}
	vec, err := loader(ctx)
	if err != nil {
		return nil, err
	}
	c.Set(ctx, key, dim, vec)
	return vec, nil
}

func (c *localVectorCache) Close() error { return c.local.Close() }

func (c *localVectorCache) observeCacheResult(result string) {
	if c.metrics != nil {
		c.metrics.ObserveCacheResult(result)
	}
}

func encodeTestVector(vec []float32) []byte {
	data := make([]byte, 4*len(vec))
	for i, v := range vec {
		binary.LittleEndian.PutUint32(data[i*4:], math.Float32bits(v))
	}
	return data
}

func decodeTestVector(data []byte, dim int) ([]float32, bool) {
	if len(data) != 4*dim {
		return nil, false
	}
	vec := make([]float32, dim)
	for i := range vec {
		vec[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[i*4:]))
	}
	return vec, true
}

func TestFdV2CacheMissThenHit(t *testing.T) {
	metrics := observability.NewMetrics()
	vectorCache := newLocalVectorCache(100, metrics)
	t.Cleanup(func() { _ = vectorCache.Close() })
	var calls atomic.Int64
	embedder := &lifecycleTestEmbedder{embedFunc: func(_ context.Context, texts []string) ([][]float32, error) {
		calls.Add(1)
		embeddings := make([][]float32, len(texts))
		for i := range texts {
			embedding := make([]float32, 1024)
			embedding[0] = 1
			embeddings[i] = embedding
		}
		return embeddings, nil
	}}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	embedHandler := handlers.NewEmbeddingsHandler(embedder, vectorCache, lifecycleTestModel, logger)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/embeddings", middleware.ValidateEmbeddingsRequest(), embedHandler.CreateEmbedding)
	r.GET("/metrics", metrics.Handler())

	first := postCacheEmbedding(t, r)
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d, want %d; body=%s", first.Code, http.StatusOK, first.Body.String())
	}
	if got := first.Header().Get(handlers.HeaderCache); got != "MISS" {
		t.Fatalf("first X-Cache = %q, want MISS", got)
	}

	started := time.Now()
	second := postCacheEmbedding(t, r)
	latency := time.Since(started)
	if second.Code != http.StatusOK {
		t.Fatalf("second status = %d, want %d; body=%s", second.Code, http.StatusOK, second.Body.String())
	}
	if got := second.Header().Get(handlers.HeaderCache); got != "HIT" {
		t.Fatalf("second X-Cache = %q, want HIT", got)
	}
	if latency >= 5*time.Millisecond {
		t.Fatalf("cache HIT latency = %s, want < 5ms", latency)
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("embedder calls = %d, want 1", got)
	}

	metricsResponse := httptest.NewRecorder()
	r.ServeHTTP(metricsResponse, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))
	if !strings.Contains(metricsResponse.Body.String(), `fd_cache_hits_total{result="hit"} 1`) {
		t.Fatalf("metrics missing hit counter:\n%s", metricsResponse.Body.String())
	}
}

func postCacheEmbedding(t *testing.T, r http.Handler) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", strings.NewReader(`{"model":"test","input":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}
