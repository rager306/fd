package main

import (
	"context"
	"errors"
	"fd-api/lifecycle"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRouterDoesNotTrustForwardedForByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	configureTrustedProxies(r)
	r.GET("/client-ip", func(c *gin.Context) { c.String(http.StatusOK, c.ClientIP()) })

	req := httptest.NewRequest(http.MethodGet, "/client-ip", http.NoBody)
	req.RemoteAddr = "198.51.100.10:1234"
	req.Header.Set("X-Forwarded-For", "203.0.113.99")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Body.String(); got != "198.51.100.10" {
		t.Fatalf("ClientIP = %q, want direct remote address", got)
	}
}

type warmupModelFunc func(ctx context.Context, texts []string) ([][]float32, error)

func (f warmupModelFunc) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return f(ctx, texts)
}

func TestStartModelWarmupMarksStateReady(t *testing.T) {
	state := lifecycle.NewState()
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		return [][]float32{{1}}, nil
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmup(logger, state, model, time.Second)
	waitForCondition(t, state.IsReady)
}

func TestStartModelWarmupStoresError(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		return nil, boom
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmup(logger, state, model, time.Second)
	waitForCondition(t, func() bool {
		return errors.Is(state.LastError(), boom)
	})
	if state.IsReady() {
		t.Fatal("state should not be ready after warmup failure")
	}
}

func TestStartModelWarmupRetriesAndMarksReady(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	var attempts atomic.Int32
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		if attempts.Add(1) == 1 {
			return nil, boom
		}
		return [][]float32{{1}}, nil
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmupWithPolicy(logger, state, model, time.Second, warmupRetryPolicy{
		maxAttempts: 3,
		backoff:     func(int) time.Duration { return 0 },
	})

	waitForCondition(t, state.IsReady)
	if attempts.Load() != 2 {
		t.Fatalf("attempts = %d, want 2", attempts.Load())
	}
	if err := state.LastError(); err != nil {
		t.Fatalf("LastError after successful retry = %v, want nil", err)
	}
}

func TestStartModelWarmupRecordsTerminalErrorAfterMaxAttempts(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	var attempts atomic.Int32
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		attempts.Add(1)
		return nil, boom
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmupWithPolicy(logger, state, model, time.Second, warmupRetryPolicy{
		maxAttempts: 3,
		backoff:     func(int) time.Duration { return 0 },
	})

	waitForCondition(t, func() bool { return attempts.Load() == 3 })
	if !errors.Is(state.LastError(), boom) {
		t.Fatalf("LastError = %v, want boom", state.LastError())
	}
	if state.IsReady() {
		t.Fatal("state should not be ready after terminal warmup failure")
	}
}

func TestReportHTTPServerErrorIgnoresWrappedServerClosed(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	shutdownCh := make(chan os.Signal, 1)

	reportHTTPServerError(logger, "127.0.0.1:0", func() error {
		return fmt.Errorf("wrapped: %w", http.ErrServerClosed)
	}, shutdownCh)

	select {
	case sig := <-shutdownCh:
		t.Fatalf("unexpected shutdown signal for ErrServerClosed: %v", sig)
	default:
	}
}

func TestReportHTTPServerErrorSignalsFatalError(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	shutdownCh := make(chan os.Signal, 1)
	boom := errors.New("accept failed")

	reportHTTPServerError(logger, "127.0.0.1:0", func() error { return boom }, shutdownCh)

	select {
	case sig := <-shutdownCh:
		serverSig, ok := sig.(serverErrorSignal)
		if !ok {
			t.Fatalf("signal type = %T, want serverErrorSignal", sig)
		}
		if !errors.Is(serverSig.err, boom) {
			t.Fatalf("signal error = %v, want %v", serverSig.err, boom)
		}
		if got := sig.String(); got != "server_error" {
			t.Fatalf("signal string = %q, want server_error", got)
		}
	default:
		t.Fatal("expected fatal listener error to trigger shutdown signal")
	}
}

func waitForCondition(t *testing.T, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("condition not met within 1s")
}

func TestLoadEmbeddingRuntimeConfigDefaultsToTEI(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", "/tmp/stale-manifest.json")

	config, err := loadEmbeddingRuntimeConfig()
	if err != nil {
		t.Fatalf("loadEmbeddingRuntimeConfig default returned error: %v", err)
	}
	if config.Backend != embeddingBackendTEI {
		t.Fatalf("backend = %q, want %q", config.Backend, embeddingBackendTEI)
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsInvalidBackend(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "candle")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected invalid backend error")
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsONNXBackend(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", "/tmp/stale-manifest.json")
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", "/tmp/tokenizer.json")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected ONNX backend to be disabled")
	}
	if got := err.Error(); !strings.Contains(got, "TEI only") && !strings.Contains(got, "supports TEI only") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmbeddingRuntimeConfigHealthReturnsSafeTEIMetadata(t *testing.T) {
	config := &embeddingRuntimeConfig{Backend: embeddingBackendTEI}
	health := config.Health("deepvk/USER-bge-m3", "v2")
	if health == nil {
		t.Fatal("TEI health metadata should not be nil")
	}
	if health.Backend != "tei" {
		t.Fatalf("backend = %q, want tei", health.Backend)
	}
	if health.Model != "deepvk/USER-bge-m3" {
		t.Fatalf("model = %q", health.Model)
	}
	if health.Dimensions != 1024 {
		t.Fatalf("dimensions = %d, want 1024", health.Dimensions)
	}
	if !health.ProductionDefault {
		t.Fatal("production_default should be true for TEI")
	}
	if health.CacheNamespace != "v2" {
		t.Fatalf("cache_namespace = %q", health.CacheNamespace)
	}
}
