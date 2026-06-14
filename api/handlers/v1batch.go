package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"log/slog"
)

const (
	maxV1BatchGroups = 100
	maxV1BatchInputs = 32
)

var errEmptyV1BatchEmbeddingResponse = errors.New("embedding backend returned no vectors")

// V1BatchHandler serves the OpenAI-compatible fd extension POST /v1/batch.
type V1BatchHandler struct {
	embedder Embedder
	cache    EmbeddingCache
	logger   *slog.Logger
}

// NewV1BatchHandler wires the embedder/cache/logger used by /v1/batch.
func NewV1BatchHandler(embedder Embedder, cache EmbeddingCache, logger *slog.Logger) *V1BatchHandler {
	return &V1BatchHandler{embedder: embedder, cache: cache, logger: logger}
}

type v1BatchRequest struct {
	Batches [][]string `json:"batches"`
}

type v1BatchResponse struct {
	Batches [][][]float32 `json:"batches"`
}

// CreateBatch serves POST /v1/batch.
func (h *V1BatchHandler) CreateBatch(c *gin.Context) {
	var req v1BatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid v1 batch request", "error", err)
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			WriteError(c, CodePayloadTooLarge, "", "request body exceeds max "+strconv.FormatInt(maxBytesErr.Limit, 10)+" bytes")
			return
		}
		WriteError(c, CodeInvalidJSON, "", "invalid JSON: "+err.Error())
		return
	}
	if !validateV1BatchRequest(c, req.Batches) {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	out := make([][][]float32, len(req.Batches))
	for batchIndex, batch := range req.Batches {
		vectors, ok := h.embedBatch(ctx, c, batchIndex, batch)
		if !ok {
			return
		}
		out[batchIndex] = vectors
	}
	c.JSON(http.StatusOK, v1BatchResponse{Batches: out})
}

func validateV1BatchRequest(c *gin.Context, batches [][]string) bool {
	if len(batches) == 0 {
		WriteError(c, CodeInputRequired, "batches", "batches is required (non-empty array of string arrays)")
		return false
	}
	if len(batches) > maxV1BatchGroups {
		WriteError(c, CodeBatchTooLarge, "batches", "batch group count "+strconv.Itoa(len(batches))+" exceeds max "+strconv.Itoa(maxV1BatchGroups))
		return false
	}
	for i, batch := range batches {
		if len(batch) == 0 {
			WriteError(c, CodeInputRequired, "batches["+strconv.Itoa(i)+"]", "inner batch must contain at least one string")
			return false
		}
		if len(batch) > maxV1BatchInputs {
			WriteError(c, CodeBatchTooLarge, "batches["+strconv.Itoa(i)+"]", "inner batch size "+strconv.Itoa(len(batch))+" exceeds max "+strconv.Itoa(maxV1BatchInputs))
			return false
		}
		for j, text := range batch {
			if len(text) > maxBatchInputChars {
				WriteError(c, CodeInputTooLong, "batches["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"]", "batches["+strconv.Itoa(i)+"]["+strconv.Itoa(j)+"] exceeds max length "+strconv.Itoa(maxBatchInputChars)+" chars")
				return false
			}
		}
	}
	return true
}

func (h *V1BatchHandler) embedBatch(ctx context.Context, c *gin.Context, batchIndex int, batch []string) ([][]float32, bool) {
	vectors := make([][]float32, len(batch))
	for textIndex, text := range batch {
		vector, err := h.cache.GetOrLoad(ctx, text, 1024, func(ctx context.Context) ([]float32, error) {
			vectors, err := h.embedder.Embed(ctx, []string{text})
			if err != nil {
				return nil, err
			}
			if len(vectors) == 0 {
				return nil, errEmptyV1BatchEmbeddingResponse
			}
			return vectors[0], nil
		})
		if err != nil {
			h.logger.Error("v1 batch embedding error", "error", err, "batch", batchIndex, "index", textIndex)
			WriteError(c, CodeInternalError, "", "embedding generation failed")
			return nil, false
		}
		vectors[textIndex] = vector
	}
	return vectors, true
}
