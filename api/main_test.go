package main

import (
	"context"
	"errors"
	"fd-api/lifecycle"
	"io"
	"log/slog"
	"strings"
	"testing"
	"time"
)

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
	waitForCondition(t, time.Second, state.IsReady)
}

func TestStartModelWarmupStoresError(t *testing.T) {
	state := lifecycle.NewState()
	boom := errors.New("boom")
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		return nil, boom
	})
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	startModelWarmup(logger, state, model, time.Second)
	waitForCondition(t, time.Second, func() bool {
		return errors.Is(state.LastError(), boom)
	})
	if state.IsReady() {
		t.Fatal("state should not be ready after warmup failure")
	}
}

func waitForCondition(t *testing.T, timeout time.Duration, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("condition not met within %s", timeout)
}

func TestGetEnvIntReturnsDefaultWhenUnset(t *testing.T) {
	t.Setenv("FD_TEST_INT", "")

	got := getEnvInt("FD_TEST_INT", 50)
	if got != 50 {
		t.Fatalf("getEnvInt unset = %d, want 50", got)
	}
}

func TestGetEnvIntParsesPositiveInteger(t *testing.T) {
	t.Setenv("FD_TEST_INT", "75")

	got := getEnvInt("FD_TEST_INT", 50)
	if got != 75 {
		t.Fatalf("getEnvInt = %d, want 75", got)
	}
}

func TestGetEnvIntReturnsDefaultForInvalidValue(t *testing.T) {
	t.Setenv("FD_TEST_INT", "12x")

	got := getEnvInt("FD_TEST_INT", 50)
	if got != 50 {
		t.Fatalf("getEnvInt invalid = %d, want 50", got)
	}
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
	// ONNX-only fields must be nil (omitted from JSON)
	if health.ArtifactVerified != nil {
		t.Fatal("artifact_verified should be nil for TEI")
	}
	if health.TokenizerVerified != nil {
		t.Fatal("tokenizer_verified should be nil for TEI")
	}
	if health.RuntimeLibraryVerified != nil {
		t.Fatal("runtime_library_verified should be nil for TEI")
	}
	if health.Provider != "" {
		t.Fatalf("provider = %q, want empty for TEI", health.Provider)
	}
}
