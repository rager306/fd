package embed

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeInner struct {
	calls      atomic.Int64
	gatherLock sync.Mutex // serializes Embed calls so coalescing isn't accidentally racy in tests
	lastBatch  atomic.Int64
}

func (f *fakeInner) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	f.calls.Add(1)
	f.gatherLock.Lock()
	defer f.gatherLock.Unlock()
	f.lastBatch.Store(int64(len(texts)))
	// Minimal work — synchronous per-call sleep to simulate real TEI latency.
	time.Sleep(2 * time.Millisecond)
	out := make([][]float32, len(texts))
	for i := range texts {
		out[i] = []float32{float32(i)}
	}
	return out, nil
}

func TestCoalescingEmbedderBurst(t *testing.T) {
	inner := &fakeInner{}
	co := NewCoalescingEmbedder(inner, 5*time.Millisecond)
	defer co.Close()

	const N = 50
	var wg sync.WaitGroup
	wg.Add(N)
	start := make(chan struct{})

	for i := 0; i < N; i++ {
		go func(i int) { //nolint:unparam // i is unused
			defer wg.Done()
			<-start
			texts := []string{
				"alpha",
				"beta",
				"gamma",
			}
			got, err := co.Embed(context.Background(), texts)
			if err != nil {
				t.Errorf("embed err: %v", err)
				return
			}
			if len(got) != 3 {
				t.Errorf("got %d embeddings, want 3", len(got))
			}
		}(i)
	}
	close(start)
	wg.Wait()

	// All 50 goroutines x 3 texts = 150 inputs. With coalescing on,
	// far fewer downstream calls than N should happen.
	calls := inner.calls.Load()
	if calls >= int64(N) {
		t.Fatalf("coalescing did not happen: %d calls for %d goroutines", calls, N)
	}
	if calls < 1 {
		t.Fatalf("expected at least 1 downstream call, got %d", calls)
	}
	t.Logf("coalesced %d concurrent jobs into %d downstream calls (%d inputs each batch)",
		N, calls, inner.lastBatch.Load())
}

func TestCoalescingEmbedderPassThroughOnZeroWindow(t *testing.T) {
	inner := &fakeInner{}
	co := NewCoalescingEmbedder(inner, 0)
	defer co.Close()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = co.Embed(context.Background(), []string{"x"})
		}()
	}
	wg.Wait()
	if calls := inner.calls.Load(); calls != 5 {
		t.Fatalf("pass-through calls = %d, want 5 (one per goroutine)", calls)
	}
}

func TestCoalescingEmbedderContextCancel(t *testing.T) {
	inner := &fakeInner{}
	co := NewCoalescingEmbedder(inner, 50*time.Millisecond) // wide window
	defer co.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := co.Embed(ctx, []string{"hi"})
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
}
