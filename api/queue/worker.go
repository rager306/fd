package queue

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"fd-api/embed"
)

// Default worker parameters. Tests may construct StartQueueWorker with
// different values via the WorkerConfig struct.
const (
	defaultBatchMaxSize = 32
	defaultBatchWindow = 10 * time.Millisecond
)

// WorkerConfig holds tunable parameters for the queue worker. Zero
// values fall back to defaults via Normalize.
type WorkerConfig struct {
	BatchMaxSize int
	BatchWindow  time.Duration
}

// Normalize fills unset fields with defaults and clamps invalid values.
func (c *WorkerConfig) Normalize() {
	if c.BatchMaxSize <= 0 {
		c.BatchMaxSize = defaultBatchMaxSize
	}
	if c.BatchWindow <= 0 {
		c.BatchWindow = defaultBatchWindow
	}
}

// StartQueueWorker drains the bounded channel and processes items as
// time-windowed batches: up to cfg.BatchMaxSize items per TEI call, or
// until cfg.BatchWindow elapses since the first item arrived in the
// current batch. Splits TEI's per-batch response back into per-item
// results stored in the result store. The worker stops when ctx is
// cancelled or the channel is closed.
func StartQueueWorker(
	ctx context.Context,
	store *ResultStore,
	items <-chan Item,
	emb embed.Embedder,
	logger *slog.Logger,
	cfg WorkerConfig,
) {
	cfg.Normalize()
	go runQueueWorker(ctx, store, items, emb, logger, cfg)
}

func runQueueWorker(
	ctx context.Context,
	store *ResultStore,
	items <-chan Item,
	emb embed.Embedder,
	logger *slog.Logger,
	cfg WorkerConfig,
) {
	logger.Info("queue worker started",
		"batch_max_size", cfg.BatchMaxSize,
		"batch_window_ms", cfg.BatchWindow.Milliseconds(),
	)
	var processed int64
	defer func() {
		logger.Info("queue worker stopped", "processed", processed)
	}()

	for {
		select {
		case <-ctx.Done():
			drainMark := drainRemaining(items, store, context.Canceled, logger)
			if drainMark > 0 {
				logger.Info("queue worker drained remaining items on shutdown", "count", drainMark)
			}
			return
		default:
		}

		// Blocking receive of first item in a batch.
		firstItem, ok := awaitFirstItem(ctx, items)
		if !ok {
			return // ctx cancelled or channel closed
		}
		batch := []Item{firstItem}
		windowStart := time.Now()
		batch = drainUpToMax(ctx, items, batch, cfg.BatchMaxSize, windowStart, cfg.BatchWindow)
		processed += int64(len(batch))

		started := time.Now()
		results := processBatch(ctx, batch, emb)
		elapsed := time.Since(started)
		for i, r := range results {
			r.CompletedAt = time.Now().UnixNano()
			r.CreatedAt = batch[i].CreatedAt
			store.Save(batch[i].ID, &r)
			batch[i].Response <- r
		}
		if batchSize := len(batch); batchSize > 1 {
			logger.Info("queue batch processed",
				"size", batchSize,
				"latency_ms", elapsed.Milliseconds(),
			)
		}
	}
}

// awaitFirstItem blocks until a new item arrives, ctx is cancelled, or
// the channel is closed. Returns (item, true) on success, (_, false) when
// the worker should exit.
func awaitFirstItem(ctx context.Context, items <-chan Item) (Item, bool) {
	select {
	case <-ctx.Done():
		return Item{}, false
	case item, ok := <-items:
		if !ok {
			return Item{}, false
		}
		return item, true
	}
}

// drainUpToMax collects more items from items into batch up to maxSize
// or until window has elapsed since windowStart. Returns the final batch.
// Honors ctx: returns current batch if cancellation is observed between
// attempts.
func drainUpToMax(ctx context.Context, items <-chan Item, batch []Item, maxSize int, windowStart time.Time, window time.Duration) []Item {
	deadline := windowStart.Add(window)
	timer := time.NewTimer(time.Until(deadline))
	defer timer.Stop()
	for len(batch) < maxSize {
		select {
		case <-ctx.Done():
			return batch
		case <-timer.C:
			return batch
		case item, ok := <-items:
			if !ok {
				return batch
			}
			batch = append(batch, item)
		default:
			// No more items immediately available. Return current batch;
			// worker will await another first item on next iteration.
			return batch
		}
	}
	return batch
}

// processBatch issues a single TEI call for the entire batch and splits
// the per-item results. On TEI failure all items are marked failed with
// the same error.
func processBatch(ctx context.Context, batch []Item, emb embed.Embedder) []Result {
	results := make([]Result, len(batch))

	// Concat all texts preserving per-item boundaries.
	var texts []string
	indexByID := make([]*Item, 0, len(batch))
	for i := range batch {
		texts = append(texts, batch[i].Texts...)
		indexByID = append(indexByID, &batch[i]) //nolint:staticcheck // deliberate
	}

	if ctx.Err() != nil {
		for i := range batch {
			results[i] = Result{
				ID:     batch[i].ID,
				Status: StatusFailed,
				Err:    ctx.Err(),
			}
		}
		return results
	}

	embs, err := emb.Embed(ctx, texts)
	if err != nil {
		for i := range batch {
			results[i] = Result{
				ID:     batch[i].ID,
				Status: StatusFailed,
				Err:    fmt.Errorf("batch TEI call failed: %w", err),
			}
		}
		return results
	}

	// Walk embeddings, slicing into per-item slices.
	cursor := 0
	for i, item := range batch {
		n := len(item.Texts)
		itemSlice := make([][]float32, n)
		if cursor+n <= len(embs) {
			for j := 0; j < n; j++ {
				itemSlice[j] = embs[cursor+j]
			}
		}
		results[i] = Result{
			ID:         item.ID,
			Status:     StatusCompleted,
			Embeddings: itemSlice,
		}
		cursor += n
	}
	return results
}

// drainRemaining marks all queued items as failed with the given context
// error. Used during shutdown so pollers don't block on items that will
// never complete.
func drainRemaining(items <-chan Item, store *ResultStore, err error, logger *slog.Logger) int {
	count := 0
	for {
		select {
		case item, ok := <-items:
			if !ok {
				return count
			}
			result := Result{
				ID:        item.ID,
				Status:    StatusFailed,
				Err:       err,
				CompletedAt: time.Now().UnixNano(),
			}
			store.Save(item.ID, &result)
			// item.Response уже закрыт через channel cancel — неблокирующий send
			select {
			case item.Response <- result:
			default:
				logger.Warn("queue drain response channel closed", "id", item.ID)
			}
			count++
		default:
			return count
		}
	}
}

