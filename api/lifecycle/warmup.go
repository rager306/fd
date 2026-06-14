package lifecycle

import (
	"context"
	"fmt"
)

const warmupInput = "fd lifecycle warmup"

// WarmupModel is the minimal model surface required by PreWarm.
type WarmupModel interface {
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}

// PreWarm issues a single dummy embedding request to force model/runtime loading.
// It validates that exactly one embedding is returned so startup cannot mark
// readiness on a malformed model response.
func PreWarm(ctx context.Context, model WarmupModel) error {
	if model == nil {
		return fmt.Errorf("warmup model is nil")
	}

	embeddings, err := model.Embed(ctx, []string{warmupInput})
	if err != nil {
		return fmt.Errorf("prewarm embedding: %w", err)
	}
	if len(embeddings) != 1 {
		return fmt.Errorf("prewarm embedding count=%d want 1", len(embeddings))
	}
	if len(embeddings[0]) == 0 {
		return fmt.Errorf("prewarm embedding vector is empty")
	}
	return nil
}
