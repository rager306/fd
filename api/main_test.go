package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fd-api/embed"
	"fd-api/lifecycle"
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
	t.Setenv("ONNX_ARTIFACT_MANIFEST", "")

	config, err := loadEmbeddingRuntimeConfig()
	if err != nil {
		t.Fatalf("loadEmbeddingRuntimeConfig default returned error: %v", err)
	}
	if config.Backend != embeddingBackendTEI {
		t.Fatalf("backend = %q, want %q", config.Backend, embeddingBackendTEI)
	}
	if config.ONNXArtifact != nil {
		t.Fatal("default TEI config should not validate ONNX artifact")
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsInvalidBackend(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "candle")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected invalid backend error")
	}
}

func TestLoadEmbeddingRuntimeConfigRequiresONNXManifest(t *testing.T) {
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", "")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected missing ONNX_ARTIFACT_MANIFEST error")
	}
}

func TestLoadEmbeddingRuntimeConfigRequiresONNXRuntimeLibrary(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, false)
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected missing ONNX_RUNTIME_LIBRARY error")
	}
}

func TestLoadEmbeddingRuntimeConfigRequiresONNXTokenizerPath(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, false)
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", "")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected missing ONNX_TOKENIZER_PATH error")
	}
}

func TestLoadEmbeddingRuntimeConfigValidatesONNXManifest(t *testing.T) {
	manifestPath, artifactPath := writeMainTestONNXManifest(t, false)
	tokenizerPath := filepath.Join(filepath.Dir(manifestPath), "tokenizer.json")
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", tokenizerPath)
	t.Setenv("ONNX_MAX_SEQUENCE_LENGTH", "128")

	config, err := loadEmbeddingRuntimeConfig()
	if err != nil {
		t.Fatalf("loadEmbeddingRuntimeConfig onnx returned error: %v", err)
	}
	if config.Backend != embeddingBackendONNX {
		t.Fatalf("backend = %q, want %q", config.Backend, embeddingBackendONNX)
	}
	if config.ONNXManifestPath != manifestPath {
		t.Fatalf("manifest path = %q, want %q", config.ONNXManifestPath, manifestPath)
	}
	if config.ONNXArtifact == nil || config.ONNXArtifact.Path != artifactPath {
		t.Fatalf("validated artifact path = %#v, want %q", config.ONNXArtifact, artifactPath)
	}
	if config.ONNXRuntimeLibraryPath != "/tmp/libonnxruntime.so" {
		t.Fatalf("runtime library path = %q", config.ONNXRuntimeLibraryPath)
	}
	if config.ONNXTokenizerPath != tokenizerPath {
		t.Fatalf("tokenizer path = %q", config.ONNXTokenizerPath)
	}
	if config.ONNXMaxSequenceLength != 128 {
		t.Fatalf("max sequence length = %d, want 128", config.ONNXMaxSequenceLength)
	}
	if config.ONNXArtifact.ValidatedMaxSequenceLength != 1024 {
		t.Fatalf("validated max sequence length = %d, want 1024", config.ONNXArtifact.ValidatedMaxSequenceLength)
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsSequenceLengthAboveManifestContract(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, false)
	tokenizerPath := filepath.Join(filepath.Dir(manifestPath), "tokenizer.json")
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", tokenizerPath)
	t.Setenv("ONNX_MAX_SEQUENCE_LENGTH", "2048")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected sequence length contract error")
	}
	if got := err.Error(); !strings.Contains(got, "ONNX_MAX_SEQUENCE_LENGTH=2048") || !strings.Contains(got, "validated_max_sequence_length=1024") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEmbeddingRuntimeConfigHealthForONNX(t *testing.T) {
	config := &embeddingRuntimeConfig{
		Backend:                    embeddingBackendONNX,
		ONNXMaxSequenceLength:      1024,
		ONNXProvider:               onnxProviderCPU,
		ONNXTokenizerVerified:      true,
		ONNXRuntimeLibraryVerified: true,
		ONNXArtifact: &embed.ONNXArtifactValidation{
			ArtifactID:                 "artifact-test",
			Dimensions:                 1024,
			ProductionDefault:          false,
			ValidatedMaxSequenceLength: 1024,
		},
	}

	health := config.Health("deepvk/USER-bge-m3", "m026-test")
	if health == nil {
		t.Fatal("expected ONNX runtime health")
	}
	if health.ArtifactID != "artifact-test" || health.CacheNamespace != "m026-test" || health.ArtifactVerified == nil || !*health.ArtifactVerified {
		t.Fatalf("unexpected health metadata: %#v", health)
	}
	if health.Provider != onnxProviderCPU || health.TokenizerVerified == nil || !*health.TokenizerVerified || health.RuntimeLibraryVerified == nil || !*health.RuntimeLibraryVerified {
		t.Fatalf("unexpected verification metadata: %#v", health)
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

func TestLoadEmbeddingRuntimeConfigRejectsTokenizerChecksumMismatch(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, false)
	tokenizerPath := filepath.Join(filepath.Dir(manifestPath), "tokenizer.json")
	if err := os.WriteFile(tokenizerPath, []byte("tampered tokenizer"), 0o600); err != nil {
		t.Fatalf("tamper tokenizer: %v", err)
	}
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", tokenizerPath)

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected tokenizer checksum error")
	}
	if got := err.Error(); !strings.Contains(got, "ONNX tokenizer JSON") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsUnsupportedProvider(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, false)
	tokenizerPath := filepath.Join(filepath.Dir(manifestPath), "tokenizer.json")
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", tokenizerPath)
	t.Setenv("ONNX_PROVIDER", "CUDAExecutionProvider")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected unsupported provider error")
	}
	if got := err.Error(); !strings.Contains(got, "ONNX_PROVIDER") || !strings.Contains(got, "CPUExecutionProvider") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadEmbeddingRuntimeConfigValidatesRuntimeLibrarySHAWhenConfigured(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, false)
	tokenizerPath := filepath.Join(filepath.Dir(manifestPath), "tokenizer.json")
	runtimeLibraryPath := filepath.Join(filepath.Dir(manifestPath), "libonnxruntime.so")
	runtimeBytes := []byte("fake onnx runtime")
	if err := os.WriteFile(runtimeLibraryPath, runtimeBytes, 0o600); err != nil {
		t.Fatalf("write runtime library: %v", err)
	}
	digest := sha256.Sum256(runtimeBytes)
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", runtimeLibraryPath)
	t.Setenv("ONNX_TOKENIZER_PATH", tokenizerPath)
	t.Setenv("ONNX_RUNTIME_SHA256", hex.EncodeToString(digest[:]))

	config, err := loadEmbeddingRuntimeConfig()
	if err != nil {
		t.Fatalf("loadEmbeddingRuntimeConfig onnx returned error: %v", err)
	}
	if !config.ONNXRuntimeLibraryVerified {
		t.Fatal("expected runtime library verification flag")
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsRuntimeLibrarySHAMismatch(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, false)
	tokenizerPath := filepath.Join(filepath.Dir(manifestPath), "tokenizer.json")
	runtimeLibraryPath := filepath.Join(filepath.Dir(manifestPath), "libonnxruntime.so")
	if err := os.WriteFile(runtimeLibraryPath, []byte("fake onnx runtime"), 0o600); err != nil {
		t.Fatalf("write runtime library: %v", err)
	}
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", runtimeLibraryPath)
	t.Setenv("ONNX_TOKENIZER_PATH", tokenizerPath)
	t.Setenv("ONNX_RUNTIME_SHA256", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected runtime library checksum error")
	}
	if got := err.Error(); !strings.Contains(got, "ONNX runtime library verification failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsInvalidONNXManifest(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, true)
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", filepath.Join(filepath.Dir(manifestPath), "tokenizer.json"))

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected invalid ONNX manifest error")
	}
}

func writeMainTestONNXManifest(t *testing.T, corruptDigest bool) (manifestPath, artifactPath string) {
	t.Helper()

	dir := t.TempDir()
	artifactRelPath := filepath.Join(".gsd", "runtime", "onnx", "m010-s03", "model.onnx")
	artifactPath = filepath.Join(dir, artifactRelPath)
	artifactBytes := []byte("fake onnx bytes for main config test")
	if err := os.MkdirAll(filepath.Dir(artifactPath), 0o700); err != nil {
		t.Fatalf("create artifact dir: %v", err)
	}
	if err := os.WriteFile(artifactPath, artifactBytes, 0o600); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
	tokenizerPath := filepath.Join(dir, "tokenizer.json")
	tokenizerBytes := []byte("fake tokenizer json for main config test")
	if err := os.WriteFile(tokenizerPath, tokenizerBytes, 0o600); err != nil {
		t.Fatalf("write tokenizer: %v", err)
	}
	tokenizerDigest := sha256.Sum256(tokenizerBytes)
	tokenizerDigestHex := hex.EncodeToString(tokenizerDigest[:])
	digest := sha256.Sum256(artifactBytes)
	digestHex := hex.EncodeToString(digest[:])
	if corruptDigest {
		digestHex = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	}

	manifest := map[string]any{
		"schema_version":     1,
		"artifact_id":        "main-test-onnx",
		"status":             "prototype_only",
		"production_default": false,
		"artifact": map[string]any{
			"local_path":  artifactRelPath,
			"size_bytes":  len(artifactBytes),
			"sha256":      digestHex,
			"git_tracked": false,
		},
		"model": map[string]any{
			"source_files": map[string]any{
				"tokenizer.json": map[string]any{
					"size_bytes": len(tokenizerBytes),
					"sha256":     tokenizerDigestHex,
				},
			},
		},
		"runtime": map[string]any{
			"outputs": []map[string]any{
				{
					"name":  "dense_vecs",
					"shape": []any{"batch_size", 1024},
					"type":  "tensor(float)",
				},
			},
			"expected_dimensions":           1024,
			"expected_normalized":           true,
			"validated_max_sequence_length": 1024,
		},
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal manifest: %v", err)
	}
	manifestPath = filepath.Join(dir, "manifest.json")
	if err := os.WriteFile(manifestPath, manifestBytes, 0o600); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	return manifestPath, artifactPath
}
