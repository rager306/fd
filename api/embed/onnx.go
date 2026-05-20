//go:build onnx

package embed

import (
	"context"
	"fmt"
	"os"
	"sync"

	ort "github.com/yalue/onnxruntime_go"
)

const defaultONNXMaxSequenceLength = 512

type ONNXEmbedder struct {
	artifact          *ONNXArtifactValidation
	tokenizer         onnxTokenizer
	session           *ort.DynamicAdvancedSession
	maxSequenceLength int
	mu                sync.Mutex
}

var onnxRuntimeEnvMu sync.Mutex

func NewONNXEmbedder(options ONNXEmbedderOptions) (*ONNXEmbedder, error) {
	if options.ManifestPath == "" {
		return nil, fmt.Errorf("onnx manifest path is required")
	}
	if options.SharedLibraryPath == "" {
		return nil, fmt.Errorf("onnx runtime shared library path is required")
	}
	if _, err := os.Stat(options.SharedLibraryPath); err != nil {
		return nil, fmt.Errorf("onnx runtime shared library invalid %q: %w", options.SharedLibraryPath, err)
	}
	if options.TokenizerPath == "" {
		return nil, fmt.Errorf("onnx tokenizer path is required")
	}
	if _, err := os.Stat(options.TokenizerPath); err != nil {
		return nil, fmt.Errorf("onnx tokenizer path invalid %q: %w", options.TokenizerPath, err)
	}
	maxSequenceLength := options.MaxSequenceLength
	if maxSequenceLength <= 0 {
		maxSequenceLength = defaultONNXMaxSequenceLength
	}

	artifact, err := ValidateONNXArtifactManifest(options.ManifestPath)
	if err != nil {
		return nil, err
	}

	tk, err := newONNXTokenizer(options.TokenizerPath)
	if err != nil {
		return nil, err
	}

	onnxRuntimeEnvMu.Lock()
	ort.SetSharedLibraryPath(options.SharedLibraryPath)
	if !ort.IsInitialized() {
		if err := ort.InitializeEnvironment(); err != nil {
			onnxRuntimeEnvMu.Unlock()
			return nil, fmt.Errorf("initialize onnxruntime environment: %w", err)
		}
	}
	onnxRuntimeEnvMu.Unlock()

	session, err := ort.NewDynamicAdvancedSession(
		artifact.Path,
		[]string{"input_ids", "attention_mask"},
		[]string{artifact.OutputName},
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create onnx session artifact_id=%q path=%q: %w", artifact.ArtifactID, artifact.Path, err)
	}

	return &ONNXEmbedder{
		artifact:          artifact,
		tokenizer:         tk,
		session:           session,
		maxSequenceLength: maxSequenceLength,
	}, nil
}

func (e *ONNXEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	inputIDs, attentionMask, sequenceLength, err := e.encodeBatch(texts)
	if err != nil {
		return nil, err
	}

	inputTensor, err := ort.NewTensor(ort.NewShape(int64(len(texts)), int64(sequenceLength)), inputIDs)
	if err != nil {
		return nil, fmt.Errorf("create onnx input_ids tensor: %w", err)
	}
	defer func() {
		_ = inputTensor.Destroy()
	}()

	attentionTensor, err := ort.NewTensor(ort.NewShape(int64(len(texts)), int64(sequenceLength)), attentionMask)
	if err != nil {
		return nil, fmt.Errorf("create onnx attention_mask tensor: %w", err)
	}
	defer func() {
		_ = attentionTensor.Destroy()
	}()

	outputs := []ort.Value{nil}
	e.mu.Lock()
	err = e.session.Run([]ort.Value{inputTensor, attentionTensor}, outputs)
	e.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("run onnx dense inference artifact_id=%q: %w", e.artifact.ArtifactID, err)
	}
	if outputs[0] == nil {
		return nil, fmt.Errorf("onnx dense inference returned nil output artifact_id=%q", e.artifact.ArtifactID)
	}
	defer func() {
		_ = outputs[0].Destroy()
	}()

	outputTensor, ok := outputs[0].(*ort.Tensor[float32])
	if !ok {
		return nil, fmt.Errorf("onnx dense output has unexpected tensor type %T", outputs[0])
	}
	data := outputTensor.GetData()
	expectedValues := len(texts) * e.artifact.Dimensions
	if len(data) != expectedValues {
		return nil, fmt.Errorf("onnx dense output size=%d want %d", len(data), expectedValues)
	}

	result := make([][]float32, len(texts))
	for i := range texts {
		start := i * e.artifact.Dimensions
		end := start + e.artifact.Dimensions
		vector := make([]float32, e.artifact.Dimensions)
		copy(vector, data[start:end])
		result[i] = vector
	}
	return result, nil
}

func (e *ONNXEmbedder) Close() error {
	if e == nil {
		return nil
	}
	var firstErr error
	if e.session != nil {
		firstErr = e.session.Destroy()
	}
	if e.tokenizer != nil {
		e.tokenizer.Close()
	}
	return firstErr
}

func (e *ONNXEmbedder) Artifact() *ONNXArtifactValidation {
	if e == nil || e.artifact == nil {
		return nil
	}
	copyValue := *e.artifact
	return &copyValue
}

func (e *ONNXEmbedder) encodeBatch(texts []string) ([]int64, []int64, int, error) {
	encoded := make([]onnxTokenEncoding, len(texts))
	sequenceLength := 0
	for i, text := range texts {
		item, err := e.tokenizer.Encode(text)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("tokenize input index=%d: %w", i, err)
		}
		if len(item.InputIDs) == 0 {
			return nil, nil, 0, fmt.Errorf("tokenize input index=%d produced no tokens", i)
		}
		if len(item.InputIDs) != len(item.AttentionMask) {
			return nil, nil, 0, fmt.Errorf("tokenize input index=%d produced ids/mask length mismatch", i)
		}
		encoded[i] = item
		length := len(item.InputIDs)
		if length > e.maxSequenceLength {
			length = e.maxSequenceLength
		}
		if length > sequenceLength {
			sequenceLength = length
		}
	}
	if sequenceLength == 0 {
		return nil, nil, 0, fmt.Errorf("onnx tokenized batch produced zero sequence length")
	}

	inputIDs := make([]int64, len(texts)*sequenceLength)
	attentionMask := make([]int64, len(texts)*sequenceLength)
	for row, item := range encoded {
		limit := len(item.InputIDs)
		if limit > sequenceLength {
			limit = sequenceLength
		}
		for col := 0; col < limit; col++ {
			offset := row*sequenceLength + col
			inputIDs[offset] = int64(item.InputIDs[col])
			attentionMask[offset] = int64(item.AttentionMask[col])
		}
	}
	return inputIDs, attentionMask, sequenceLength, nil
}
