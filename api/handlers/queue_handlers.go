package handlers

import (
	"net/http"
	"strings"
	"time"

	"fd-api/embed"
	"fd-api/queue"

	"github.com/gin-gonic/gin"
	"log/slog"
)

// queueRetryAfterSeconds is the Retry-After hint sent on queue_full
// submissions. Picked conservative so clients back off predictably.
const queueRetryAfterSeconds = "5"

// QueueHandler serves POST /v1/queue (submit) and GET /v1/queue/:id (poll).
// The handler is constructed only when FD_QUEUE_ENABLED=true; otherwise
// callers see 404 not_found from the standard NoRoute handler.
type QueueHandler struct {
	store     *queue.ResultStore
	items     chan<- queue.Item
	logger    *slog.Logger
	modelID   string
}

// NewQueueHandler wires the bounded channel, result store, and model identity.
func NewQueueHandler(store *queue.ResultStore, items chan<- queue.Item, modelID string, logger *slog.Logger) *QueueHandler {
	return &QueueHandler{store: store, items: items, modelID: modelID, logger: logger}
}

// Submit handles POST /v1/queue. Validates the OpenAI-compatible body and
// enqueues the inputs. On queue saturation returns 503 with a Retry-After
// hint so producers back off under burst load.
func (h *QueueHandler) Submit(c *gin.Context) {
	var req embed.EmbeddingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, CodeInvalidJSON, "", "invalid JSON: "+err.Error())
		return
	}
	if len(req.Input) == 0 {
		WriteError(c, CodeInputRequired, "input", "input is required (non-empty array of strings)")
		return
	}
	if len(req.Input) > 128 {
		WriteError(c, CodeBatchTooLarge, "input", "batch size exceeds max 128; split into smaller batches")
		return
	}
	for _, t := range req.Input {
		if strings.TrimSpace(t) == "" {
			WriteError(c, CodeInputRequired, "input", "input contains empty strings")
			return
		}
	}
	dims := 1024
	if req.Dimensions != nil && *req.Dimensions != 0 {
		dims = *req.Dimensions
	}

	id := queue.NewRequestID()
	item := queue.Item{
		ID:        id,
		Texts:     req.Input,
		Dims:      dims,
		Response:  make(chan queue.Result, 1),
		SubmitCtx: c.Request.Context(),
		CreatedAt: time.Now().UnixNano(),
	}

	// Non-blocking send with saturation backpressure.
	select {
	case h.items <- item:
		// Accepted. Track pending result so GET /v1/queue/:id can see it
		// even before the worker has filled in the response.
		h.store.Save(id, &queue.Result{
			ID:        id,
			Status:    queue.StatusPending,
			CreatedAt: item.CreatedAt,
		})
		c.Header("X-Request-Id", id)
		c.JSON(http.StatusAccepted, gin.H{
			"status":     queue.StatusPending,
			"request_id": id,
			"hint":       "Poll GET /v1/queue/" + id,
		})
	default:
		c.Header("Retry-After", queueRetryAfterSeconds)
		WriteError(c, CodeModelOverloaded, "", "queue is saturated; retry shortly")
	}
}

// Poll handles GET /v1/queue/:id. Returns pending (202), completed (200
// with embedding data), failed (500 with error envelope), or 404 when
// the id is unknown or has expired from the result store.
func (h *QueueHandler) Poll(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		WriteError(c, CodeInputRequired, "id", "missing request id")
		return
	}
	res, ok := h.store.Get(id)
	if !ok {
		WriteError(c, CodeNotFound, "id", "queue request not found")
		return
	}
	c.Header("X-Request-Id", res.ID)
	switch res.Status {
	case queue.StatusPending:
		c.JSON(http.StatusAccepted, gin.H{
			"status":     queue.StatusPending,
			"request_id": res.ID,
		})
	case queue.StatusCompleted:
		data := make([]embed.EmbeddingObj, len(res.Embeddings))
		for i, emb := range res.Embeddings {
			obj := embed.EmbeddingObj{
				Object:    "embedding",
				Index:     i,
				Dimensions: 1024,
			}
			obj.SetVector(emb)
			data[i] = obj
		}
		c.JSON(http.StatusOK, embed.EmbeddingsResponse{
			Object: "list", //nolint:goconst // literal used for json response
			Data:   data,
			Model:  h.modelID,
			Usage: embed.Usage{
				PromptTokens: res.PromptTokens,
				TotalTokens:  res.PromptTokens,
			},
		})
	case queue.StatusFailed:
		WriteError(c, CodeInternalError, "", "queue processing failed: "+safeErr(res.Err))
	default:
		WriteError(c, CodeInternalError, "", "unknown queue status: "+string(res.Status))
	}
}

// safeErr returns the error string or a placeholder when the worker wrote
// a nil error on a failed entry (defensive).
func safeErr(err error) string {
	if err == nil {
		return "unknown"
	}
	return err.Error()
}
