//go:build !onnx

package embed

import (
	"context"
	"fmt"
)

type ONNXEmbedder struct{}

func NewONNXEmbedder(options ONNXEmbedderOptions) (*ONNXEmbedder, error) {
	return nil, fmt.Errorf("onnx backend requires build tag onnx")
}

func (e *ONNXEmbedder) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	return nil, fmt.Errorf("onnx backend requires build tag onnx")
}

func (e *ONNXEmbedder) Close() error {
	return nil
}
