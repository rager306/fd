package lifecycle

import (
	"context"
	"errors"
	"testing"
)

type warmupModelFunc func(ctx context.Context, texts []string) ([][]float32, error)

func (f warmupModelFunc) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	return f(ctx, texts)
}

func TestPreWarmCallsModelWithSingleDummyInput(t *testing.T) {
	var gotTexts []string
	model := warmupModelFunc(func(_ context.Context, texts []string) ([][]float32, error) {
		gotTexts = append([]string(nil), texts...)
		return [][]float32{{1, 2, 3}}, nil
	})

	if err := PreWarm(context.Background(), model); err != nil {
		t.Fatalf("PreWarm() = %v, want nil", err)
	}
	if len(gotTexts) != 1 {
		t.Fatalf("warmup input count = %d, want 1", len(gotTexts))
	}
	if gotTexts[0] != warmupInput {
		t.Fatalf("warmup input = %q, want %q", gotTexts[0], warmupInput)
	}
}

func TestPreWarmWrapsModelError(t *testing.T) {
	boom := errors.New("boom")
	model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
		return nil, boom
	})

	err := PreWarm(context.Background(), model)
	if !errors.Is(err, boom) {
		t.Fatalf("PreWarm() = %v, want wrapping boom", err)
	}
}

func TestPreWarmRejectsMalformedResponses(t *testing.T) {
	tests := []struct {
		name       string
		embeddings [][]float32
	}{
		{name: "zero embeddings", embeddings: nil},
		{name: "two embeddings", embeddings: [][]float32{{1}, {2}}},
		{name: "empty vector", embeddings: [][]float32{{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := warmupModelFunc(func(_ context.Context, _ []string) ([][]float32, error) {
				return tt.embeddings, nil
			})

			if err := PreWarm(context.Background(), model); err == nil {
				t.Fatal("PreWarm() = nil, want malformed response error")
			}
		})
	}
}

func TestPreWarmPropagatesContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	model := warmupModelFunc(func(ctx context.Context, _ []string) ([][]float32, error) {
		return nil, ctx.Err()
	})

	err := PreWarm(ctx, model)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("PreWarm() = %v, want context canceled", err)
	}
}

func TestPreWarmRejectsNilModel(t *testing.T) {
	if err := PreWarm(context.Background(), nil); err == nil {
		t.Fatal("PreWarm(nil) = nil, want error")
	}
}
