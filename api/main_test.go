package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

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
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", "/tmp/tokenizer.json")
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
	if config.ONNXTokenizerPath != "/tmp/tokenizer.json" {
		t.Fatalf("tokenizer path = %q", config.ONNXTokenizerPath)
	}
	if config.ONNXMaxSequenceLength != 128 {
		t.Fatalf("max sequence length = %d, want 128", config.ONNXMaxSequenceLength)
	}
}

func TestLoadEmbeddingRuntimeConfigRejectsInvalidONNXManifest(t *testing.T) {
	manifestPath, _ := writeMainTestONNXManifest(t, true)
	t.Setenv("EMBEDDING_BACKEND", "onnx")
	t.Setenv("ONNX_ARTIFACT_MANIFEST", manifestPath)
	t.Setenv("ONNX_RUNTIME_LIBRARY", "/tmp/libonnxruntime.so")
	t.Setenv("ONNX_TOKENIZER_PATH", "/tmp/tokenizer.json")

	_, err := loadEmbeddingRuntimeConfig()
	if err == nil {
		t.Fatal("expected invalid ONNX manifest error")
	}
}

func writeMainTestONNXManifest(t *testing.T, corruptDigest bool) (manifestPath string, artifactPath string) {
	t.Helper()

	dir := t.TempDir()
	artifactPath = filepath.Join(dir, "model.onnx")
	artifactBytes := []byte("fake onnx bytes for main config test")
	if err := os.WriteFile(artifactPath, artifactBytes, 0o600); err != nil {
		t.Fatalf("write artifact: %v", err)
	}
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
			"local_path":  artifactPath,
			"size_bytes":  len(artifactBytes),
			"sha256":      digestHex,
			"git_tracked": false,
		},
		"runtime": map[string]any{
			"outputs": []map[string]any{
				{
					"name":  "dense_vecs",
					"shape": []any{"batch_size", 1024},
					"type":  "tensor(float)",
				},
			},
			"expected_dimensions": 1024,
			"expected_normalized": true,
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
