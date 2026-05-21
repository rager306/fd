package embed

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	ONNXArtifactStatusPrototypeOnly = "prototype_only"
	ONNXExpectedOutputName          = "dense_vecs"
	ONNXExpectedDimensions          = 1024
)

var (
	ErrONNXManifestMissingArtifact   = errors.New("onnx artifact missing")
	ErrONNXManifestChecksumMismatch  = errors.New("onnx artifact checksum mismatch")
	ErrONNXManifestMetadataMismatch  = errors.New("onnx artifact metadata mismatch")
	ErrONNXManifestProductionDefault = errors.New("onnx artifact manifest must not be production default")
)

type ONNXArtifactManifest struct {
	SchemaVersion     int    `json:"schema_version"`
	ArtifactID        string `json:"artifact_id"`
	Status            string `json:"status"`
	ProductionDefault bool   `json:"production_default"`
	manifestDir       string
	Model             struct {
		SourceFiles struct {
			TokenizerJSON struct {
				SizeBytes int64  `json:"size_bytes"`
				SHA256    string `json:"sha256"`
			} `json:"tokenizer.json"`
		} `json:"source_files"`
	} `json:"model"`
	Artifact struct {
		LocalPath  string `json:"local_path"`
		SizeBytes  int64  `json:"size_bytes"`
		SHA256     string `json:"sha256"`
		GitTracked bool   `json:"git_tracked"`
	} `json:"artifact"`
	Runtime struct {
		Outputs []struct {
			Name  string        `json:"name"`
			Shape []interface{} `json:"shape"`
			Type  string        `json:"type"`
		} `json:"outputs"`
		ExpectedDimensions         int  `json:"expected_dimensions"`
		ExpectedNormalized         bool `json:"expected_normalized"`
		ValidatedMaxSequenceLength int  `json:"validated_max_sequence_length"`
	} `json:"runtime"`
}

type ONNXArtifactValidation struct {
	ArtifactID                 string
	Path                       string
	SizeBytes                  int64
	SHA256                     string
	OutputName                 string
	Dimensions                 int
	ProductionDefault          bool
	ValidatedMaxSequenceLength int
	TokenizerJSONSizeBytes     int64
	TokenizerJSONSHA256        string
}

func LoadONNXArtifactManifest(path string) (*ONNXArtifactManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read onnx artifact manifest %q: %w", path, err)
	}

	var manifest ONNXArtifactManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("parse onnx artifact manifest %q: %w", path, err)
	}
	manifest.manifestDir = filepath.Dir(path)
	return &manifest, nil
}

func ValidateONNXArtifactManifest(path string) (*ONNXArtifactValidation, error) {
	manifest, err := LoadONNXArtifactManifest(path)
	if err != nil {
		return nil, err
	}
	return manifest.ValidateArtifact()
}

func (m *ONNXArtifactManifest) ValidateArtifact() (*ONNXArtifactValidation, error) {
	if m.ProductionDefault {
		return nil, fmt.Errorf("%w: artifact_id=%q", ErrONNXManifestProductionDefault, m.ArtifactID)
	}
	if m.Artifact.GitTracked {
		return nil, fmt.Errorf("%w: artifact_id=%q git_tracked=true", ErrONNXManifestMetadataMismatch, m.ArtifactID)
	}
	if m.Artifact.LocalPath == "" {
		return nil, fmt.Errorf("%w: artifact_id=%q missing artifact.local_path", ErrONNXManifestMetadataMismatch, m.ArtifactID)
	}
	if err := validateONNXArtifactLocalPath(m.Artifact.LocalPath); err != nil {
		return nil, fmt.Errorf("%w: artifact_id=%q artifact.local_path=%q: %w", ErrONNXManifestMetadataMismatch, m.ArtifactID, safePathDisplay(m.Artifact.LocalPath), err)
	}
	if m.Artifact.SizeBytes <= 0 {
		return nil, fmt.Errorf("%w: artifact_id=%q invalid artifact.size_bytes=%d", ErrONNXManifestMetadataMismatch, m.ArtifactID, m.Artifact.SizeBytes)
	}
	if m.Model.SourceFiles.TokenizerJSON.SizeBytes < 0 {
		return nil, fmt.Errorf("%w: artifact_id=%q tokenizer.json size_bytes=%d", ErrONNXManifestMetadataMismatch, m.ArtifactID, m.Model.SourceFiles.TokenizerJSON.SizeBytes)
	}
	if m.Model.SourceFiles.TokenizerJSON.SHA256 != "" && len(m.Model.SourceFiles.TokenizerJSON.SHA256) != sha256HexLength {
		return nil, fmt.Errorf("%w: artifact_id=%q invalid tokenizer.json sha256", ErrONNXManifestMetadataMismatch, m.ArtifactID)
	}
	if m.Runtime.ExpectedDimensions != ONNXExpectedDimensions {
		return nil, fmt.Errorf("%w: artifact_id=%q expected_dimensions=%d want %d", ErrONNXManifestMetadataMismatch, m.ArtifactID, m.Runtime.ExpectedDimensions, ONNXExpectedDimensions)
	}
	if !m.Runtime.ExpectedNormalized {
		return nil, fmt.Errorf("%w: artifact_id=%q expected_normalized=false", ErrONNXManifestMetadataMismatch, m.ArtifactID)
	}
	if m.Runtime.ValidatedMaxSequenceLength < 0 {
		return nil, fmt.Errorf("%w: artifact_id=%q validated_max_sequence_length=%d", ErrONNXManifestMetadataMismatch, m.ArtifactID, m.Runtime.ValidatedMaxSequenceLength)
	}
	if len(m.Runtime.Outputs) != 1 {
		return nil, fmt.Errorf("%w: artifact_id=%q outputs=%d want 1", ErrONNXManifestMetadataMismatch, m.ArtifactID, len(m.Runtime.Outputs))
	}
	output := m.Runtime.Outputs[0]
	if output.Name != ONNXExpectedOutputName {
		return nil, fmt.Errorf("%w: artifact_id=%q output=%q want %q", ErrONNXManifestMetadataMismatch, m.ArtifactID, output.Name, ONNXExpectedOutputName)
	}
	if !outputShapeHasExpectedDimensions(output.Shape) {
		return nil, fmt.Errorf("%w: artifact_id=%q output shape does not include dimension %d", ErrONNXManifestMetadataMismatch, m.ArtifactID, ONNXExpectedDimensions)
	}

	artifactPath := resolveONNXArtifactPath(m.manifestDir, m.Artifact.LocalPath)
	artifactDisplay := safePathDisplay(m.Artifact.LocalPath)
	info, err := os.Stat(artifactPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: artifact_id=%q path=%q", ErrONNXManifestMissingArtifact, m.ArtifactID, artifactDisplay)
		}
		return nil, fmt.Errorf("stat onnx artifact %q: %w", artifactDisplay, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%w: artifact_id=%q path=%q is directory", ErrONNXManifestMissingArtifact, m.ArtifactID, artifactDisplay)
	}
	if info.Size() != m.Artifact.SizeBytes {
		return nil, fmt.Errorf("%w: artifact_id=%q path=%q size=%d want %d", ErrONNXManifestMetadataMismatch, m.ArtifactID, artifactDisplay, info.Size(), m.Artifact.SizeBytes)
	}

	digest, err := sha256File(artifactPath)
	if err != nil {
		return nil, err
	}
	if digest != m.Artifact.SHA256 {
		return nil, fmt.Errorf("%w: artifact_id=%q path=%q expected=%s actual=%s", ErrONNXManifestChecksumMismatch, m.ArtifactID, artifactDisplay, m.Artifact.SHA256, digest)
	}

	return &ONNXArtifactValidation{
		ArtifactID:                 m.ArtifactID,
		Path:                       artifactPath,
		SizeBytes:                  info.Size(),
		SHA256:                     digest,
		OutputName:                 output.Name,
		Dimensions:                 m.Runtime.ExpectedDimensions,
		ProductionDefault:          m.ProductionDefault,
		ValidatedMaxSequenceLength: m.Runtime.ValidatedMaxSequenceLength,
		TokenizerJSONSizeBytes:     m.Model.SourceFiles.TokenizerJSON.SizeBytes,
		TokenizerJSONSHA256:        m.Model.SourceFiles.TokenizerJSON.SHA256,
	}, nil
}

const sha256HexLength = 64

var allowedONNXArtifactPathPrefixes = []string{
	".gsd/runtime/onnx/",
	".gsd/runtime/tokenizers/",
	".gsd/runtime/onnxruntime/",
	"tei-models/",
}

func validateONNXArtifactLocalPath(path string) error {
	if filepath.IsAbs(path) {
		return errors.New("absolute paths are not allowed")
	}
	cleaned := filepath.ToSlash(filepath.Clean(path))
	if cleaned == "." || cleaned == ".." || cleaned == "" || strings.HasPrefix(cleaned, "../") {
		return errors.New("path must be repo-relative and must not traverse outside the repository")
	}
	for _, prefix := range allowedONNXArtifactPathPrefixes {
		if strings.TrimSuffix(prefix, "/") == cleaned || strings.HasPrefix(cleaned, prefix) {
			return nil
		}
	}
	return fmt.Errorf("path must be under one of the approved artifact roots: %v", allowedONNXArtifactPathPrefixes)
}

func safePathDisplay(path string) string {
	if path == "" {
		return ""
	}
	cleaned := filepath.ToSlash(filepath.Clean(path))
	if filepath.IsAbs(path) {
		return ".../" + filepath.Base(path)
	}
	return cleaned
}

func resolveONNXArtifactPath(manifestDir, artifactPath string) string {
	if filepath.IsAbs(artifactPath) {
		return artifactPath
	}
	if _, err := os.Stat(artifactPath); err == nil {
		return artifactPath
	}
	for dir := manifestDir; dir != "." && dir != string(filepath.Separator); dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, artifactPath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		next := filepath.Dir(dir)
		if next == dir {
			break
		}
	}
	return artifactPath
}

func outputShapeHasExpectedDimensions(shape []interface{}) bool {
	for _, value := range shape {
		switch v := value.(type) {
		case float64:
			if int(v) == ONNXExpectedDimensions {
				return true
			}
		case int:
			if v == ONNXExpectedDimensions {
				return true
			}
		}
	}
	return false
}

func sha256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open onnx artifact %q: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("hash onnx artifact %q: %w", path, err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
