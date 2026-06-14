package main

import (
	"context"
	"io"
	"log/slog"
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

func TestFdV2CacheMissThenHit(t *testing.T) {
	metrics := observability.NewMetrics()
	vectorCache := cache.NewLRUCache(100, time.Hour, metrics)
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
