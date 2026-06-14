// Package handlers implements the legacy /embeddings/batch endpoint used by FalkorDB callers. Preserves base64-by-default response shape for backward compatibility.
package handlers

import (
	"context"
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
	teiClient Embedder
	cache     EmbeddingCache
	modelID   string
	logger    *slog.Logger
}

// NewBatchHandler wires the embedder, cache, model ID, and logger used by
// the legacy batch embeddings endpoint.
func NewBatchHandler(teiClient Embedder, c EmbeddingCache, modelID string, logger *slog.Logger) *BatchHandler {
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
		WriteError(c, CodeInvalidJSON, "", "invalid JSON: "+err.Error())
		return
	}

	if len(req.Inputs) == 0 {
		WriteError(c, CodeInputRequired, "inputs", "inputs is required (non-empty array of strings)")
		return
	}

	dims := 1024
	if req.Dimensions != 0 {
		if req.Dimensions != 512 && req.Dimensions != 1024 {
			WriteError(c, CodeDimensionsInvalid, "dimensions",
				"dimensions must be 1024 or 512, got "+strconv.Itoa(req.Dimensions))
			return
		}
		dims = req.Dimensions
	}

	if req.EncodingFormat != "" && req.EncodingFormat != embed.EncodingFormatBase64 && req.EncodingFormat != embed.EncodingFormatFloat {
		WriteError(c, CodeEncodingInvalid, "encoding_format",
			"encoding_format must be float or base64, got \""+req.EncodingFormat+"\"")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	embeddings := make([]string, len(req.Inputs))
	for i, text := range req.Inputs {
		emb, err := h.cache.GetOrLoad(ctx, text, dims, func(ctx context.Context) ([]float32, error) {
			embs, err := h.teiClient.Embed(ctx, []string{text})
			if err != nil {
				return nil, err
			}
			return embs[0], nil
		})
		if err != nil {
			h.logger.Error("batch embedding error", "error", err, "index", i)
			WriteError(c, CodeInternalError, "", "embedding generation failed")
			return
		}

		fullEmb := emb
		if dims == 512 && len(fullEmb) >= 512 {
			fullEmb = fullEmb[:512]
		}

		embeddings[i] = embed.EncodeEmbedding(fullEmb, defaultBatchEncoding(req.EncodingFormat))
	}

	c.JSON(http.StatusOK, embed.BatchEmbeddingsResponse{
		Embeddings: embeddings,
		Count:      len(req.Inputs),
		Dimensions: dims,
	})
}

// (Removed local itoa — strconv.Itoa is sufficient.)

// defaultBatchEncoding preserves the legacy /embeddings/batch default of
// base64 (FalkorDB callers depend on it). The OpenAI-style /v1/embeddings
// endpoint defaults to float instead.
func defaultBatchEncoding(format string) string {
	if format == "" {
		return embed.EncodingFormatBase64
	}
	return format
}
