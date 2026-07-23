package embed

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// coalesceBatchFull is the hard cap per TEI batch sent to the downstream
// embedder. Aligned with the TEI max_client_batch_size of 32.
const coalesceBatchFull = 32

// coalescedJob represents one caller waiting for an embedding result.
type coalescedJob struct {
	texts  []string
	result chan coalescedResult
}

type coalescedResult struct {
	embeddings [][]float32
	err        error
}

// CoalescingEmbedder wraps an underlying Embedder and time-windowed batching.
// Concurrent calls to Embed within coalesceWindow are merged into a single
// downstream call (up to coalesceBatchFull inputs), reducing TEI call count
// under burst. Concurrent callers block until the batched call returns;
// the window is short enough (default 5 ms) to be transparent for most
// same-host consumers.
type CoalescingEmbedder struct {
	inner  Embedder
	jobs   chan coalescedJob
	window time.Duration
	wg     sync.WaitGroup
}

// NewCoalescingEmbedder starts a background batcher and returns a wrapper.
// window is the per-batch collection interval. Set to 0 to disable coalescing
// (Embed calls pass through directly to inner).
func NewCoalescingEmbedder(inner Embedder, window time.Duration) *CoalescingEmbedder {
	if inner == nil {
		return nil
	}
	c := &CoalescingEmbedder{
		inner: inner,
		jobs:   make(chan coalescedJob),
		window: window,
	}
	if window > 0 {
		c.wg.Add(1)
		go c.run()
	}
	return c
}

// Close shuts down the background batcher. Pending jobs are drained before
// calling the downstream embedder; this may block for up to window time.
func (c *CoalescingEmbedder) Close() {
	close(c.jobs)
	c.wg.Wait()
}

// Embed satisfies the Embedder interface. When coalescing is active (window >0)
// the call is enqueued into the background batcher. When window is zero the
// call falls through directly to the inner embedder.
func (c *CoalescingEmbedder) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	if c.window <= 0 {
		return c.inner.Embed(ctx, texts)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	result := make(chan coalescedResult, 1)
	select {
	case c.jobs <- coalescedJob{texts: texts, result: result}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	select {
	case r := <-result:
		return r.embeddings, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *CoalescingEmbedder) run() {
	defer c.wg.Done()
	var batch []coalescedJob
	for {
		job, ok := <-c.jobs
		if !ok {
			// Shutdown: drain whatever is left.
			if len(batch) > 0 {
				c.flushBatch(batch)
			}
			return
		}
		batch = append(batch, job)
		// Try to fill batch up to coalesceBatchFull without blocking.
		timer := time.NewTimer(c.window)
		for len(batch) < coalesceBatchFull {
			select {
			case next, open := <-c.jobs:
				if !open {
					c.flushBatch(batch)
					return
				}
				batch = append(batch, next)
			case <-timer.C:
				// Window expired — flush what we have.
				c.flushBatch(batch)
				batch = batch[:0]
				goto outer
			}
		}
		// Batch full — flush immediately.
		c.flushBatch(batch)
		batch = batch[:0]
	outer:
	}
}

func (c *CoalescingEmbedder) flushBatch(batch []coalescedJob) {
	if len(batch) == 0 {
		return
	}
	// Collect all texts preserving per-job boundaries.
	var allTexts []string
	counts := make([]int, len(batch))
	for i, j := range batch {
		counts[i] = len(j.texts)
		allTexts = append(allTexts, j.texts...)
	}
	embs, err := c.inner.Embed(context.Background(), allTexts)
	// Split by job.
	cursor := 0
	for i, j := range batch {
		n := counts[i]
		if err != nil {//nolint:gocritic // readability
			j.result <- coalescedResult{err: err}
		} else if cursor+n <= len(embs) {
			slice := make([][]float32, n)
			copy(slice, embs[cursor:cursor+n])
			j.result <- coalescedResult{embeddings: slice}
		} else {
			j.result <- coalescedResult{err: fmt.Errorf("coalesce split: cursor %d+%d > len(embs) %d", cursor, n, len(embs))}
		}
		cursor += n
	}
}


