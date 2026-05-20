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
	Artifact          struct {
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
		ExpectedDimensions int  `json:"expected_dimensions"`
		ExpectedNormalized bool `json:"expected_normalized"`
	} `json:"runtime"`
}

type ONNXArtifactValidation struct {
	ArtifactID string
	Path       string
	SizeBytes  int64
	SHA256     string
	OutputName string
	Dimensions int
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
	if m.Artifact.SizeBytes <= 0 {
		return nil, fmt.Errorf("%w: artifact_id=%q invalid artifact.size_bytes=%d", ErrONNXManifestMetadataMismatch, m.ArtifactID, m.Artifact.SizeBytes)
	}
	if len(m.Artifact.SHA256) != sha256HexLength {
		return nil, fmt.Errorf("%w: artifact_id=%q invalid artifact.sha256", ErrONNXManifestMetadataMismatch, m.ArtifactID)
	}
	if m.Runtime.ExpectedDimensions != ONNXExpectedDimensions {
		return nil, fmt.Errorf("%w: artifact_id=%q expected_dimensions=%d want %d", ErrONNXManifestMetadataMismatch, m.ArtifactID, m.Runtime.ExpectedDimensions, ONNXExpectedDimensions)
	}
	if !m.Runtime.ExpectedNormalized {
		return nil, fmt.Errorf("%w: artifact_id=%q expected_normalized=false", ErrONNXManifestMetadataMismatch, m.ArtifactID)
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
	info, err := os.Stat(artifactPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: artifact_id=%q path=%q", ErrONNXManifestMissingArtifact, m.ArtifactID, m.Artifact.LocalPath)
		}
		return nil, fmt.Errorf("stat onnx artifact %q: %w", artifactPath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%w: artifact_id=%q path=%q is directory", ErrONNXManifestMissingArtifact, m.ArtifactID, artifactPath)
	}
	if info.Size() != m.Artifact.SizeBytes {
		return nil, fmt.Errorf("%w: artifact_id=%q path=%q size=%d want %d", ErrONNXManifestMetadataMismatch, m.ArtifactID, artifactPath, info.Size(), m.Artifact.SizeBytes)
	}

	digest, err := sha256File(artifactPath)
	if err != nil {
		return nil, err
	}
	if digest != m.Artifact.SHA256 {
		return nil, fmt.Errorf("%w: artifact_id=%q path=%q expected=%s actual=%s", ErrONNXManifestChecksumMismatch, m.ArtifactID, artifactPath, m.Artifact.SHA256, digest)
	}

	return &ONNXArtifactValidation{
		ArtifactID: m.ArtifactID,
		Path:       artifactPath,
		SizeBytes:  info.Size(),
		SHA256:     digest,
		OutputName: output.Name,
		Dimensions: m.Runtime.ExpectedDimensions,
	}, nil
}

const sha256HexLength = 64

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
