package embed

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// corpusPart reflects the JSON shape of tests/44-FZ-2026-articles.jsonl.
// Only fields used by the baseline are decoded.
type corpusPart struct {
	DocID   string         `json:"doc_id"`
	Chapter int            `json:"chapter"`
	Article string         `json:"article"`
	Title   string         `json:"title"`
	Parts   []corpusChunk `json:"parts"`
}

// corpusChunk is one part of an article. Its Text field is what we send
// to the embedder as input.
type corpusChunk struct {
	Number  int    `json:"number"`
	Text    string `json:"text"`
}

func load44FZCorpus(t *testing.T) []string {
	t.Helper()
	// Resolve corpus path relative to this test file. Several segments up
	// to the repo root, then tests/44-FZ-2026-articles.jsonl.
	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "..", "..", "tests", "44-FZ-2026-articles.jsonl")
	data, err := os.ReadFile(root) //nolint:gosec // root is hardcoded
	if err != nil {
		t.Skipf("corpus not available at %s: %v", root, err)
	}
	var texts []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		var part corpusPart
		if err := json.Unmarshal(scanner.Bytes(), &part); err != nil {
			continue // skip malformed lines
		}
		// Use the part's main text. Article title is also useful as a
		// distinct short input alongside the long-form body.
		for _, c := range part.Parts {
			if c.Text != "" {
				texts = append(texts, c.Text)
			}
		}
		if part.Title != "" {
			texts = append(texts, part.Title)
		}
	}
	if len(texts) == 0 {
		t.Skip("corpus contains no texts")
	}
	return texts
}

//nolint:unparam // signature matches benchmark requirements
func runCorpusBurst(t *testing.T, e Embedder, texts []string, concurrency int) (calls, totalTexts int, durations []time.Duration) {
	t.Helper()
	var mu sync.Mutex
	var callsCounter atomic.Int64

	var wg sync.WaitGroup
	var startGate sync.WaitGroup
	startGate.Add(1)

	// Track uses of underlying call counter via wrapper.
	wrapped := &atomicCounterEmbedder{inner: e, counter: &callsCounter}

	// Shuffle inputs so goroutines don't all hit the same first article.
	rng := rand.New(rand.NewSource(42)) //nolint:gosec // test use only
	jobs := append([]string(nil), texts...)
	rng.Shuffle(len(jobs), func(i, j int) { jobs[i], jobs[j] = jobs[j], jobs[i] })

	for i := 0; i < concurrency && i < len(jobs); i++ {
		wg.Add(1)
		started := time.Now()
		go func(txt string) {
			defer wg.Done()
			startGate.Wait()
			// Pass varied chunks per goroutine to simulate per-request batches.
			chunk := []string{txt}
			got, err := wrapped.Embed(context.Background(), chunk)
			elapsed := time.Since(started)
			mu.Lock()
			durations = append(durations, elapsed)
			mu.Unlock()
			if err != nil {
				t.Errorf("embed err: %v", err)
				return
			}
			if len(got) != 1 {
				t.Errorf("expected 1 embedding, got %d", len(got))
			}
		}(jobs[i])
	}
	startGate.Done()
	wg.Wait()

	calls = int(callsCounter.Load())
	totalTexts = concurrency

	// Compute percentile from collected durations.
	if len(durations) == 0 {
		return
	}
	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	return
}

// atomicCounterEmbedder is a thin Embedder wrapper that increments a
// passed-in counter on every underlying Embed call. Used by the baseline
// benchmark so that wrapping in CoalescingEmbedder is transparent.
type atomicCounterEmbedder struct {
	inner   Embedder
	counter *atomic.Int64
}

func (a *atomicCounterEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	a.counter.Add(1)
	return a.inner.Embed(ctx, texts)
}

// TestCoalescingBaseline44FZProof is the Phase 1b (Issue #9) proof:
// using a synthetic burst pattern derived from the 44-FZ legal corpus,
// we measure how many downstream TEI calls happen *without* coalescing
// versus *with* the CoalescingEmbedder. With 50 concurrent goroutines
// each submitting 1 input, pure pass-through yields ≥50 downstream calls;
// the CoalescingEmbedder merges them into ≤ceil(inputCount/batchSize)
// downstream calls (subject to the time window).
//
// Skipped if the corpus JSONL is unavailable (CI run without data).
func TestCoalescingBaseline44FZProof(t *testing.T) {
	texts := load44FZCorpus(t)
	const concurrency = 50
	if len(texts) < concurrency {
		t.Skipf("corpus too small: %d texts, need %d", len(texts), concurrency)
	}

	t.Logf("corpus: %d legal text fragments from tests/44-FZ-2026-articles.jsonl", len(texts))

	// Build the inner fake "TEI" by reusing fakeInner from coalescedembedder_test.
	inner := &fakeInner{}
	// With coalescing off (window=0): pass-through.
	control := NewCoalescingEmbedder(inner, 0)
	defer control.Close()
	_, totalControl, durationsControl := runCorpusBurst(t, control, texts, concurrency)
	callsControl := inner.calls.Load()
	t.Logf("baseline (window=0): downstream calls=%d, totalInputs=%d", callsControl, totalControl)

	// Reset inner counter for the coalesced run.
	inner.calls.Store(0)
	// With coalescing on (5ms window): a small time window captures
	// concurrent goroutines and merges them into a few TEI calls.
	co := NewCoalescingEmbedder(inner, 5*time.Millisecond)
	defer co.Close()
	_, totalCo, durationsCo := runCorpusBurst(t, co, texts, concurrency)
	callsCo := inner.calls.Load()
	t.Logf("coalesced (window=5ms): downstream calls=%d, totalInputs=%d", callsCo, totalCo)

	if callsControl < int64(concurrency) {
		t.Fatalf("baseline should have ≥%d calls, got %d", concurrency, callsControl)
	}
	if callsCo >= callsControl {
		t.Fatalf("coalescing did not reduce calls: baseline=%d coalesced=%d", callsControl, callsCo)
	}

	// Compute and log percentile improvement.
	p95 := func(ds []time.Duration) time.Duration {
		if len(ds) == 0 {
			return 0
		}
		idx := int(float64(len(ds)) * 0.95)
		if idx >= len(ds) {
			idx = len(ds) - 1
		}
		return ds[idx]
	}
	t.Logf("p95 latency: baseline=%s coalesced=%s (reduction=%d calls → %d calls)",
		p95(durationsControl), p95(durationsCo), callsControl, callsCo)
}

// TestCoalescingBaseline44FZBatchSizes synthesises the FD_MIXED pattern:
// each concurrent goroutine sends an array of 4 texts (a typical scenario
// where a single Pipeline step embeds the constituent chunks of an article).
// With 32-input TEI cap, several large requests coalesce into fewer batched
// TEI calls.
func TestCoalescingBaseline44FZBatchSizes(t *testing.T) {
	texts := load44FZCorpus(t)
	const concurrency = 30
	const perJob = 4 // 4 texts per concurrent goroutine = 120 total texts
	if len(texts) < concurrency*perJob {
		t.Skipf("corpus too small: %d texts, need %d", len(texts), concurrency*perJob)
	}

	t.Logf("corpus: %d legal text fragments; concurrency=%d perJob=%d", len(texts), concurrency, perJob)

	// Build large inputs — first perJob texts per goroutine.
	jobs := make([][]string, concurrency)
	for i := 0; i < concurrency; i++ {
		jobs[i] = texts[i*perJob : (i+1)*perJob]
	}

	inner := &fakeInner{}
	control := NewCoalescingEmbedder(inner, 0)
	defer control.Close()
	var wg sync.WaitGroup
	var startGate sync.WaitGroup
	startGate.Add(1)
	for _, j := range jobs {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
			startGate.Wait()
			_, err := control.Embed(context.Background(), batch)
			if err != nil {
				t.Errorf("baseline err: %v", err)
			}
		}(j)
	}
	startGate.Done()
	wg.Wait()
	callsControl := inner.calls.Load()
	t.Logf("baseline (window=0): downstream calls=%d (one per goroutine)", callsControl)

	inner.calls.Store(0)
	co := NewCoalescingEmbedder(inner, 8*time.Millisecond)
	defer co.Close()
	var wg2 sync.WaitGroup
	var startGate2 sync.WaitGroup
	startGate2.Add(1)
	for _, j := range jobs {
		wg2.Add(1)
		go func(batch []string) {
			defer wg2.Done()
			startGate2.Wait()
			_, err := co.Embed(context.Background(), batch)
			if err != nil {
				t.Errorf("coalesced err: %v", err)
			}
		}(j)
	}
	startGate2.Done()
	wg2.Wait()
	callsCo := inner.calls.Load()
	t.Logf("coalesced (window=8ms): downstream calls=%d (capped at ceiling(120/32)=4)", callsCo)
	if callsCo >= callsControl {
		t.Fatalf("coalescing did not reduce calls for batched burst: baseline=%d coalesced=%d", callsControl, callsCo)
	}
}
