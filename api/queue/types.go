// Package queue provides an async bulk-ingestion path for fd embedding
// requests. Submitters POST to /v1/queue and receive a request_id; a
// background worker drains the bounded queue, batches items within a
// short time-window, calls TEI once per batch, and persists per-item
// results in an in-memory store keyed by request_id. Clients poll
// GET /v1/queue/:id to fetch results.
//
// The queue is opt-in (FD_QUEUE_ENABLED, default false) and additive: the
// sync /v1/embeddings path is unchanged.
package queue

import (
	"context"
	"errors"
)

// Status is the lifecycle of a submitted queue item.
type Status string

const (
	//nolint:revive // comment not needed
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// Result holds the outcome of a worker-processed queue item. Embeddings
// has the same length as the original input. Err is nil on success.
type Result struct {
	ID           string
	Status       Status
	Embeddings   [][]float32
	PromptTokens int
	Err          error
	CreatedAt    int64 // unix nanos when submitted
	CompletedAt  int64 // unix nanos when worker finished (0 while pending)
}

// Item is the in-flight value workers consume from the bounded channel.
// Response is filled by the worker so the submitter's poll path can finish
// the request without contacting TEI directly. SubmitCtx should be the
// request's context (cancellation here stops waiting).
type Item struct {
	ID         string
	Texts      []string
	Dims       int
	Response   chan Result // buffered 1; worker writes once, submitter reads once
	SubmitCtx  context.Context
	CreatedAt  int64
}

// ErrQueueDisabled indicates the queue feature was queried while the
// FD_QUEUE_ENABLED flag is false. Returned by handlers when the route
// is unavailable.
var ErrQueueDisabled = errors.New("queue feature is disabled")

// ErrQueueFull is returned by handlers when the bounded queue rejects a
// submission because the channel buffer is at capacity.
var ErrQueueFull = errors.New("queue is at capacity")
