package embed

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateONNXArtifactManifestValid(t *testing.T) {
	manifestPath, artifactPath, digest := writeTestONNXManifest(t, func(m map[string]any) {})

	validation, err := ValidateONNXArtifactManifest(manifestPath)

	require.NoError(t, err)
	require.Equal(t, "user-bge-m3-dense-fp32-test", validation.ArtifactID)
	require.Equal(t, artifactPath, validation.Path)
	require.Equal(t, int64(len(testONNXArtifactBytes())), validation.SizeBytes)
	require.Equal(t, digest, validation.SHA256)
	require.Equal(t, ONNXExpectedOutputName, validation.OutputName)
	require.Equal(t, ONNXExpectedDimensions, validation.Dimensions)
	require.Equal(t, 1024, validation.ValidatedMaxSequenceLength)
	tokenizerDigest := sha256.Sum256(testTokenizerJSONBytes())
	require.Equal(t, int64(len(testTokenizerJSONBytes())), validation.TokenizerJSONSizeBytes)
	require.Equal(t, hex.EncodeToString(tokenizerDigest[:]), validation.TokenizerJSONSHA256)
	require.False(t, validation.ProductionDefault)
}

func TestValidateONNXArtifactManifestMissingArtifact(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		artifact := m["artifact"].(map[string]any)
		artifact["local_path"] = ".gsd/runtime/onnx/m010-s03/missing.onnx"
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestMissingArtifact)
	require.Contains(t, err.Error(), "artifact_id=\"user-bge-m3-dense-fp32-test\"")
}

func TestValidateONNXArtifactManifestChecksumMismatch(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		artifact := m["artifact"].(map[string]any)
		artifact["sha256"] = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestChecksumMismatch)
	require.Contains(t, err.Error(), "expected=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
}

func TestValidateONNXArtifactManifestInvalidOutputName(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		runtime := m["runtime"].(map[string]any)
		outputs := runtime["outputs"].([]map[string]any)
		outputs[0]["name"] = "last_hidden_state"
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestMetadataMismatch)
	require.Contains(t, err.Error(), "last_hidden_state")
}

func TestValidateONNXArtifactManifestInvalidDimensions(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		runtime := m["runtime"].(map[string]any)
		runtime["expected_dimensions"] = 512
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestMetadataMismatch)
	require.Contains(t, err.Error(), "expected_dimensions=512")
}

func TestValidateONNXArtifactManifestProductionDefaultRejected(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		m["production_default"] = true
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestProductionDefault)
}

func TestValidateONNXArtifactManifestRejectsNegativeValidatedMaxSequenceLength(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		runtime := m["runtime"].(map[string]any)
		runtime["validated_max_sequence_length"] = -1
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestMetadataMismatch)
	require.Contains(t, err.Error(), "validated_max_sequence_length=-1")
}

func TestValidateONNXArtifactManifestRejectsRepoExternalPath(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		artifact := m["artifact"].(map[string]any)
		artifact["local_path"] = "../outside/model.onnx"
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestMetadataMismatch)
	require.Contains(t, err.Error(), "repo-relative")
}

func TestValidateONNXArtifactManifestRejectsUnapprovedRoot(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		artifact := m["artifact"].(map[string]any)
		artifact["local_path"] = "tmp/model.onnx"
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestMetadataMismatch)
	require.Contains(t, err.Error(), "approved artifact roots")
}

func TestValidateONNXArtifactManifestRejectsInvalidTokenizerJSONSHA(t *testing.T) {
	manifestPath, _, _ := writeTestONNXManifest(t, func(m map[string]any) {
		model := m["model"].(map[string]any)
		sourceFiles := model["source_files"].(map[string]any)
		tokenizer := sourceFiles["tokenizer.json"].(map[string]any)
		tokenizer["sha256"] = "bad"
	})

	_, err := ValidateONNXArtifactManifest(manifestPath)

	require.ErrorIs(t, err, ErrONNXManifestMetadataMismatch)
	require.Contains(t, err.Error(), "invalid tokenizer.json sha256")
}

func TestLoadONNXArtifactManifestInvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	require.NoError(t, os.WriteFile(path, []byte("{"), 0o600))

	_, err := LoadONNXArtifactManifest(path)

	require.Error(t, err)
	require.False(t, errors.Is(err, ErrONNXManifestMetadataMismatch))
}

func writeTestONNXManifest(t *testing.T, mutate func(map[string]any)) (manifestPath, artifactPath, digest string) {
	t.Helper()

	dir := t.TempDir()
	artifactRelPath := filepath.Join(".gsd", "runtime", "onnx", "m010-s03", "user-bge-m3-dense.onnx")
	artifactPath = filepath.Join(dir, artifactRelPath)
	artifactBytes := testONNXArtifactBytes()
	require.NoError(t, os.MkdirAll(filepath.Dir(artifactPath), 0o700))
	require.NoError(t, os.WriteFile(artifactPath, artifactBytes, 0o600))
	tokenizerDigest := sha256.Sum256(testTokenizerJSONBytes())
	rawDigest := sha256.Sum256(artifactBytes)
	digest = hex.EncodeToString(rawDigest[:])

	manifest := map[string]any{
		"schema_version":     1,
		"artifact_id":        "user-bge-m3-dense-fp32-test",
		"status":             ONNXArtifactStatusPrototypeOnly,
		"production_default": false,
		"description":        "test manifest",
		"model": map[string]any{
			"id": "deepvk/USER-bge-m3",
			"source_files": map[string]any{
				"tokenizer.json": map[string]any{
					"size_bytes": len(testTokenizerJSONBytes()),
					"sha256":     hex.EncodeToString(tokenizerDigest[:]),
				},
			},
		},
		"artifact": map[string]any{
			"format":      "onnx",
			"dtype":       "fp32",
			"local_path":  artifactRelPath,
			"size_bytes":  len(artifactBytes),
			"sha256":      digest,
			"git_tracked": false,
		},
		"runtime": map[string]any{
			"provider": "CPUExecutionProvider",
			"outputs": []map[string]any{
				{
					"name":  ONNXExpectedOutputName,
					"shape": []any{"batch_size", ONNXExpectedDimensions},
					"type":  "tensor(float)",
				},
			},
			"expected_dimensions":           ONNXExpectedDimensions,
			"expected_normalized":           true,
			"validated_max_sequence_length": 1024,
		},
	}
	mutate(manifest)

	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	require.NoError(t, err)
	manifestPath = filepath.Join(dir, "manifest.json")
	require.NoError(t, os.WriteFile(manifestPath, manifestBytes, 0o600))
	return manifestPath, artifactPath, digest
}

func testONNXArtifactBytes() []byte {
	return []byte("fake onnx bytes for manifest validation")
}

func testTokenizerJSONBytes() []byte {
	return []byte("fake tokenizer json for manifest validation")
}
