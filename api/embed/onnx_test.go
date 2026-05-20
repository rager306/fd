//go:build onnx

package embed

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewONNXEmbedderRequiresManifestPath(t *testing.T) {
	_, err := NewONNXEmbedder(ONNXEmbedderOptions{})
	require.ErrorContains(t, err, "onnx manifest path is required")
}

func TestNewONNXEmbedderRequiresSharedLibraryPath(t *testing.T) {
	manifestPath := writeONNXEmbedderTestManifest(t, "model.onnx", []byte("fake"))
	_, err := NewONNXEmbedder(ONNXEmbedderOptions{ManifestPath: manifestPath})
	require.ErrorContains(t, err, "onnx runtime shared library path is required")
}

func TestNewONNXEmbedderRejectsMissingSharedLibraryPath(t *testing.T) {
	manifestPath := writeONNXEmbedderTestManifest(t, "model.onnx", []byte("fake"))
	_, err := NewONNXEmbedder(ONNXEmbedderOptions{ManifestPath: manifestPath, SharedLibraryPath: filepath.Join(t.TempDir(), "missing.so")})
	require.ErrorContains(t, err, "onnx runtime shared library invalid")
}

func TestNewONNXEmbedderRequiresTokenizerPath(t *testing.T) {
	sharedLibrary := filepath.Join(t.TempDir(), "libonnxruntime.so")
	require.NoError(t, os.WriteFile(sharedLibrary, []byte("fake so"), 0o600))
	manifestPath := writeONNXEmbedderTestManifest(t, "model.onnx", []byte("fake"))

	_, err := NewONNXEmbedder(ONNXEmbedderOptions{ManifestPath: manifestPath, SharedLibraryPath: sharedLibrary})

	require.ErrorContains(t, err, "onnx tokenizer path is required")
}

func TestONNXEmbedderLiveLocalArtifact(t *testing.T) {
	sharedLibrary := os.Getenv("FD_TEST_ONNX_RUNTIME_LIBRARY")
	if sharedLibrary == "" {
		t.Skip("set FD_TEST_ONNX_RUNTIME_LIBRARY to run live ONNX embedder test")
	}
	manifestPath := os.Getenv("FD_TEST_ONNX_ARTIFACT_MANIFEST")
	if manifestPath == "" {
		manifestPath = "../docs/onnx-artifacts/user-bge-m3-dense-fp32.json"
	}
	tokenizerPath := os.Getenv("FD_TEST_ONNX_TOKENIZER_PATH")
	if tokenizerPath == "" {
		tokenizerPath = "../tei-models/deepvk--USER-bge-m3/tokenizer.json"
	}

	embedder, err := NewONNXEmbedder(ONNXEmbedderOptions{
		ManifestPath:      manifestPath,
		SharedLibraryPath: sharedLibrary,
		TokenizerPath:     tokenizerPath,
		MaxSequenceLength: 128,
	})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, embedder.Close())
	}()

	embeddings, err := embedder.Embed(context.Background(), []string{"юридическая справка"})
	require.NoError(t, err)
	require.Len(t, embeddings, 1)
	require.Len(t, embeddings[0], ONNXExpectedDimensions)

	var norm float64
	for _, value := range embeddings[0] {
		norm += float64(value * value)
	}
	require.InDelta(t, 1.0, norm, 0.01)
}

func writeONNXEmbedderTestManifest(t *testing.T, artifactName string, artifactBytes []byte) string {
	t.Helper()

	dir := t.TempDir()
	artifactPath := filepath.Join(dir, artifactName)
	require.NoError(t, os.WriteFile(artifactPath, artifactBytes, 0o600))
	digest := sha256.Sum256(artifactBytes)

	manifest := map[string]any{
		"schema_version":     1,
		"artifact_id":        "onnx-embedder-test",
		"status":             ONNXArtifactStatusPrototypeOnly,
		"production_default": false,
		"artifact": map[string]any{
			"local_path":  artifactPath,
			"size_bytes":  len(artifactBytes),
			"sha256":      hex.EncodeToString(digest[:]),
			"git_tracked": false,
		},
		"runtime": map[string]any{
			"outputs": []map[string]any{
				{
					"name":  ONNXExpectedOutputName,
					"shape": []any{"batch_size", ONNXExpectedDimensions},
					"type":  "tensor(float)",
				},
			},
			"expected_dimensions": ONNXExpectedDimensions,
			"expected_normalized": true,
		},
	}
	manifestBytes, err := json.Marshal(manifest)
	require.NoError(t, err)
	manifestPath := filepath.Join(dir, "manifest.json")
	require.NoError(t, os.WriteFile(manifestPath, manifestBytes, 0o600))
	return manifestPath
}
