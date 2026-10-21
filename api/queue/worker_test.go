package queue

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

type fakeEmbedder struct {
	calls   atomic.Int64
	delay   time.Duration
	returns [][]float32
	err     error
	// Captures the number of texts seen in each Embed call (last value wins).
	lastBatchSize atomic.Int64
	// shortBy drops that many embeddings from the response, simulating a
	// backend that returns fewer vectors than inputs.
	shortBy int
}

func (f *fakeEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	f.calls.Add(1)
	f.lastBatchSize.Store(int64(len(texts)))
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if f.err != nil {
		return nil, f.err
	}
	out := make([][]float32, len(texts))
	for i := range texts {
		if i < len(f.returns) {
			out[i] = f.returns[i]
		} else {
			out[i] = []float32{float32(i)}
		}
	}
	if f.shortBy > 0 && f.shortBy < len(out) {
		out = out[:len(out)-f.shortBy]
	}
	return out, nil
}

func makeItem(id string, texts []string) Item {
	return Item{
		ID:        id,
		Texts:     texts,
		Dims:      1024,
		Response:  make(chan Result, 1),
		SubmitCtx: context.Background(),
		CreatedAt: time.Now().UnixNano(),
	}
}

func TestWorkerProcessesSingleItem(t *testing.T) {
	store := NewResultStore()
	defer func() { _ = store.Close() }()
	items := make(chan Item, 1)
	emb := &fakeEmbedder{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	StartQueueWorker(ctx, store, items, emb, discardLogger(), WorkerConfig{
		BatchMaxSize: 32,
		BatchWindow:  10 * time.Millisecond,
	})

	item := makeItem("test-1", []string{"hello", "world"})
	items <- item

	select {
	case res := <-item.Response:
		if res.Status != StatusCompleted {
			t.Fatalf("status = %s, want completed", res.Status)
		}
		if len(res.Embeddings) != 2 {
			t.Fatalf("embeddings = %d, want 2", len(res.Embeddings))
		}
		if emb.calls.Load() != 1 {
			t.Fatalf("embedder calls = %d, want 1", emb.calls.Load())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not complete within 2s")
	}
}

func TestWorkerHandlesEmbedError(t *testing.T) {
	store := NewResultStore()
	defer func() { _ = store.Close() }()
	items := make(chan Item, 1)
	emb := &fakeEmbedder{err: errors.New("boom")}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	StartQueueWorker(ctx, store, items, emb, discardLogger(), WorkerConfig{
		BatchMaxSize: 32,
		BatchWindow:  10 * time.Millisecond,
	})

	item := makeItem("test-err", []string{"x"})
	items <- item

	select {
	case res := <-item.Response:
		if res.Status != StatusFailed {
			t.Fatalf("status = %s, want failed", res.Status)
		}
		if res.Err == nil {
			t.Fatal("err should be populated")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not process within 2s")
	}

	got, ok := store.Get("test-err")
	if !ok {
		t.Fatal("result not in store")
	}
	if got.Status != StatusFailed {
		t.Fatalf("store status = %s, want failed", got.Status)
	}
}

func TestWorkerStopsOnContextCancel(t *testing.T) {
	store := NewResultStore()
	defer func() { _ = store.Close() }()
	items := make(chan Item, 4)
	emb := &fakeEmbedder{delay: 50 * time.Millisecond}

	ctx, cancel := context.WithCancel(context.Background())
	StartQueueWorker(ctx, store, items, emb, discardLogger(), WorkerConfig{
		BatchMaxSize: 32,
		BatchWindow:  10 * time.Millisecond,
	})

	item := makeItem("test-cancel", []string{"x"})
	items <- item
	cancel()

	// Either: worker was processing item and ctx fires mid-TEI → failed,
	// or: worker was in drain mode → also failed with ctx.Canceled.
	select {
	case res := <-item.Response:
		if res.Status != StatusFailed {
			t.Fatalf("status = %s, want failed", res.Status)
		}
	case <-time.After(2 * time.Second):
		got, ok := store.Get("test-cancel")
		if !ok || got.Status != StatusFailed {
			t.Fatal("worker did not record failed result on cancel")
		}
	}
}

// TestWorkerBatchesMultipleItems proves the time-windowed coalescing:
// submit 3 items quickly, expect 1 TEI call with 3 texts combined.
func TestWorkerBatchesMultipleItems(t *testing.T) {
	store := NewResultStore()
	defer func() { _ = store.Close() }()
	items := make(chan Item, 8)
	emb := &fakeEmbedder{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	StartQueueWorker(ctx, store, items, emb, discardLogger(), WorkerConfig{
		BatchMaxSize: 32,
		BatchWindow:  50 * time.Millisecond, // wide window to catch all 3
	})

	it1 := makeItem("b1", []string{"a"})
	it2 := makeItem("b2", []string{"b"})
	it3 := makeItem("b3", []string{"c"})
	items <- it1
	items <- it2
	items <- it3

	// Read individual response channels.
	for _, it := range []*Item{&it1, &it2, &it3} {
		select {
		case res := <-it.Response:
			if res.Status != StatusCompleted {
				t.Fatalf("status = %s, want completed", res.Status)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timeout waiting for %s", it.ID)
		}
	}

	// Core assertion: all three items were batched into ONE TEI call.
	if emb.calls.Load() != 1 {
		t.Fatalf("embedder calls = %d, want 1 (batched)", emb.calls.Load())
	}
	if got := emb.lastBatchSize.Load(); got != 3 {
		t.Fatalf("batch size = %d, want 3", got)
	}
}

// TestWorkerRespectsMaxBatchSize verifies that batch cap is enforced even
// when many items arrive in the same window.
func TestWorkerRespectsMaxBatchSize(t *testing.T) {
	store := NewResultStore()
	defer func() { _ = store.Close() }()
	items := make(chan Item, 20)
	emb := &fakeEmbedder{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	StartQueueWorker(ctx, store, items, emb, discardLogger(), WorkerConfig{
		BatchMaxSize: 4,
		BatchWindow:  50 * time.Millisecond,
	})

	// Submit 8 items. With maxBatchSize=4 and a wide window, expect 2
	// TEI calls each with 4 texts.
	for i := 0; i < 8; i++ {
		items <- makeItem(string(rune('A'+i)), []string{"x"})
	}

	time.Sleep(200 * time.Millisecond)
	if got := emb.calls.Load(); got != 2 {
		t.Fatalf("embedder calls = %d, want 2", got)
	}
}

// A backend that returns fewer vectors than inputs must fail the whole batch;
// previously the tail items were marked completed with nil vectors.
func TestWorkerFailsBatchOnResponseCountMismatch(t *testing.T) {
	store := NewResultStore()
	defer func() { _ = store.Close() }()
	items := make(chan Item, 1)
	emb := &fakeEmbedder{shortBy: 1}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	StartQueueWorker(ctx, store, items, emb, discardLogger(), WorkerConfig{
		BatchMaxSize: 32,
		BatchWindow:  10 * time.Millisecond,
	})

	item := makeItem("short-1", []string{"a", "b"})
	items <- item

	select {
	case res := <-item.Response:
		if res.Status != StatusFailed {
			t.Fatalf("status = %s, want failed", res.Status)
		}
		if res.Err == nil {
			t.Fatal("failed result carries no error")
		}
		if res.Embeddings != nil {
			t.Fatalf("failed result carries embeddings: %v", res.Embeddings)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not complete within 2s")
	}
}

// The result must report the dimensions and wire format the caller asked for,
// and account for prompt tokens instead of leaving usage at zero.
func TestWorkerHonoursRequestedDimensionsAndCountsTokens(t *testing.T) {
	store := NewResultStore()
	defer func() { _ = store.Close() }()
	items := make(chan Item, 1)
	full := make([]float32, 1024)
	for i := range full {
		full[i] = float32(i)
	}
	emb := &fakeEmbedder{returns: [][]float32{full}}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	StartQueueWorker(ctx, store, items, emb, discardLogger(), WorkerConfig{
		BatchMaxSize: 32,
		BatchWindow:  10 * time.Millisecond,
	})

	item := makeItem("dims-1", []string{"0123456789abcdefghij"})
	item.Dims = 512
	item.EncodingFormat = "base64"
	items <- item

	select {
	case res := <-item.Response:
		if res.Status != StatusCompleted {
			t.Fatalf("status = %s, want completed", res.Status)
		}
		if len(res.Embeddings) != 1 {
			t.Fatalf("embeddings = %d, want 1", len(res.Embeddings))
		}
		if got := len(res.Embeddings[0]); got != 512 {
			t.Fatalf("vector dims = %d, want 512", got)
		}
		if res.Dims != 512 {
			t.Fatalf("reported dims = %d, want 512", res.Dims)
		}
		if res.EncodingFormat != "base64" {
			t.Fatalf("encoding format = %q, want base64", res.EncodingFormat)
		}
		if res.PromptTokens == 0 {
			t.Fatal("prompt tokens = 0, want a non-zero count")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("worker did not complete within 2s")
	}
}
