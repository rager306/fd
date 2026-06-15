// Package handlers implements the legacy /embeddings/batch endpoint used by FalkorDB callers. Preserves base64-by-default response shape for backward compatibility.
package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"fd-api/embed"

	"github.com/gin-gonic/gin"
	"log/slog"
)

// BatchHandler serves the legacy /embeddings/batch endpoint used by
// FalkorDB integrations that expect base64 vectors by default.
type BatchHandler struct {
	teiClient embed.Embedder
	cache     EmbeddingCache
	modelID   string
	logger    *slog.Logger
}

// NewBatchHandler wires the embedder, cache, model ID, and logger used by
// the legacy batch embeddings endpoint.
func NewBatchHandler(teiClient embed.Embedder, c EmbeddingCache, modelID string, logger *slog.Logger) *BatchHandler {
	return &BatchHandler{
		teiClient: teiClient,
		cache:     c,
		modelID:   modelID,
		logger:    logger,
	}
}

// CreateBatchEmbeddings serves POST /embeddings/batch (legacy FalkorDB endpoint).
// /v1/embeddings uses middleware validation; this handler keeps its own
// inline validation because the request shape (inputs/dimensions) differs
// from /v1/embeddings (input/dimensions).
func (h *BatchHandler) CreateBatchEmbeddings(c *gin.Context) {
	var req embed.BatchEmbeddingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid batch request", "error", err)
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			WriteError(c, CodePayloadTooLarge, "", "request body exceeds max "+strconv.FormatInt(maxBytesErr.Limit, 10)+" bytes")
			return
		}
		WriteError(c, CodeInvalidJSON, "", "invalid JSON: "+err.Error())
		return
	}

	dims, ok := validateLegacyBatchRequest(c, req)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	vectors, err := loadBatchEmbeddings(ctx, h.cache, h.teiClient, req.Inputs, dims, batchTEISubBatchSize)
	if err != nil {
		h.logger.Error("batch embedding error", "error", err, "input_count", len(req.Inputs))
		WriteError(c, CodeInternalError, "", "embedding generation failed")
		return
	}

	embeddings := make([]string, len(vectors))
	for i, vector := range vectors {
		embeddings[i] = embed.EncodeEmbedding(vector, defaultBatchEncoding(req.EncodingFormat))
	}

	c.JSON(http.StatusOK, embed.BatchEmbeddingsResponse{
		Embeddings: embeddings,
		Count:      len(req.Inputs),
		Dimensions: dims,
	})
}

func validateLegacyBatchRequest(c *gin.Context, req embed.BatchEmbeddingsRequest) (int, bool) {
	if len(req.Inputs) == 0 {
		WriteError(c, CodeInputRequired, "inputs", "inputs is required (non-empty array of strings)")
		return 0, false
	}
	if len(req.Inputs) > maxLegacyBatchInputs {
		WriteError(c, CodeBatchTooLarge, "inputs", "batch size "+strconv.Itoa(len(req.Inputs))+" exceeds max "+strconv.Itoa(maxLegacyBatchInputs))
		return 0, false
	}
	for i, text := range req.Inputs {
		if len(text) > maxBatchInputChars {
			WriteError(c, CodeInputTooLong, "inputs["+strconv.Itoa(i)+"]", "input["+strconv.Itoa(i)+"] exceeds max length "+strconv.Itoa(maxBatchInputChars)+" chars")
			return 0, false
		}
	}

	dims := 1024
	if req.Dimensions != 0 {
		if req.Dimensions != 512 && req.Dimensions != 1024 {
			WriteError(c, CodeDimensionsInvalid, "dimensions", "dimensions must be 1024 or 512, got "+strconv.Itoa(req.Dimensions))
			return 0, false
		}
		dims = req.Dimensions
	}

	if req.EncodingFormat != "" && req.EncodingFormat != embed.EncodingFormatBase64 && req.EncodingFormat != embed.EncodingFormatFloat {
		WriteError(c, CodeEncodingInvalid, "encoding_format", "encoding_format must be float or base64, got \""+req.EncodingFormat+"\"")
		return 0, false
	}
	return dims, true
}

// defaultBatchEncoding preserves the legacy /embeddings/batch default of
// base64 (FalkorDB callers depend on it). The OpenAI-style /v1/embeddings
// endpoint defaults to float instead.
func defaultBatchEncoding(format string) string {
	if format == "" {
		return embed.EncodingFormatBase64
	}
	return format
}
