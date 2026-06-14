//go:build !onnx

package embed

import (
	"context"
	"fmt"
)

// ONNXEmbedder is a placeholder implementation compiled when the onnx
// build tag is not enabled. It makes the package API available while
// failing closed if callers try to construct or use ONNX inference.
type ONNXEmbedder struct{}

// NewONNXEmbedder returns an error when fd is built without the onnx build tag.
func NewONNXEmbedder(options ONNXEmbedderOptions) (*ONNXEmbedder, error) {
	return nil, fmt.Errorf("onnx backend requires build tag onnx")
}

// Embed returns an error when fd is built without the onnx build tag.
func (e *ONNXEmbedder) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	return nil, fmt.Errorf("onnx backend requires build tag onnx")
}

// Close is a no-op for the disabled ONNX embedder placeholder.
func (e *ONNXEmbedder) Close() error {
	return nil
}
