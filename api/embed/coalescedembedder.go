package embed

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// coalesceBatchFull is the hard cap per TEI batch sent to the downstream
// embedder. Aligned with the TEI max_client_batch_size of 32, this limit
// counts individual texts rather than jobs so a few large requests cannot
// overflow TEI's input capacity.
const coalesceBatchMaxTexts = 32

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
	// ctx bounds the downstream calls issued by the batcher. It is cancelled by
	// Close so a shutting-down process does not keep issuing TEI requests.
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCoalescingEmbedder starts a background batcher and returns a wrapper.
// window is the per-batch collection interval. Set to 0 to disable coalescing
// (Embed calls pass through directly to inner).
func NewCoalescingEmbedder(inner Embedder, window time.Duration) *CoalescingEmbedder {
	if inner == nil {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	c := &CoalescingEmbedder{
		inner:  inner,
		jobs:   make(chan coalescedJob),
		window: window,
		ctx:    ctx,
		cancel: cancel,
	}
	if window > 0 {
		c.wg.Add(1)
		go c.run()
	}
	return c
}

// Close shuts down the background batcher: it stops accepting jobs, cancels the
// in-flight downstream call, and waits for the batcher goroutine to return.
// Pending jobs fail with context.Canceled rather than being dropped silently, so
// callers do not wait on a result that will never come.
func (c *CoalescingEmbedder) Close() {
	c.cancel()
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
		// Try to fill batch up to coalesceBatchMaxTexts texts without
		// blocking. Cap is on individual texts to respect TEI's
		// max_client_batch_size=32; large jobs may exhaust the budget alone.
		timer := time.NewTimer(c.window)
		textsInBatch := totalTexts(batch)
		for textsInBatch < coalesceBatchMaxTexts {
			select {
			case next, open := <-c.jobs:
				if !open {
					c.flushBatch(batch)
					return
				}
				if textsInBatch > 0 && textsInBatch+len(next.texts) > coalesceBatchMaxTexts {
					// This job would push the batch past TEI's input cap, which
					// rejects the whole request: flush what we have and carry the
					// job into the next batch instead.
					c.flushBatch(batch)
					batch = append(batch[:0], next)
					textsInBatch = len(next.texts)
					continue
				}
				batch = append(batch, next)
				textsInBatch += len(next.texts)
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

// totalTexts returns the count of individual inputs across all jobs in
// the batch. Used to cap a coalesced batch at TEI's max input capacity.
func totalTexts(batch []coalescedJob) int {
	sum := 0
	for _, j := range batch {
		sum += len(j.texts)
	}
	return sum
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
	embs, err := c.inner.Embed(c.ctx, allTexts)
	// Split by job.
	cursor := 0
	for i, j := range batch {
		n := counts[i]
		switch {
		case err != nil:
			j.result <- coalescedResult{err: err}
		case cursor+n <= len(embs):
			slice := make([][]float32, n)
			copy(slice, embs[cursor:cursor+n])
			j.result <- coalescedResult{embeddings: slice}
		default:
			j.result <- coalescedResult{err: fmt.Errorf("coalesce split: cursor %d+%d > len(embs) %d", cursor, n, len(embs))}
		}
		cursor += n
	}
}
