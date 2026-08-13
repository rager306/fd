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
	defer store.Close()
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
	defer store.Close()
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
