package handlers

import (
	"context"
	"net/http"
	"time"

	"fd-api/embed"
	"fd-api/cache"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type EmbeddingsHandler struct {
	teiClient *embed.TEIClient
	cache     *cache.RedisCache
	modelID   string
	logger    *slog.Logger
}

func NewEmbeddingsHandler(teiClient *embed.TEIClient, cache *cache.RedisCache, modelID string, logger *slog.Logger) *EmbeddingsHandler {
	return &EmbeddingsHandler{
		teiClient: teiClient,
		cache:     cache,
		modelID:   modelID,
		logger:    logger,
	}
}

func (h *EmbeddingsHandler) CreateEmbedding(c *gin.Context) {
	var req embed.EmbeddingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Normalize: if input is a single string, convert to slice
	texts := req.Input
	if len(texts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "input is required"})
		return
	}

	// Dimensions: nil=1024 (default), 512, or explicit 1024
	dims := 1024
	if req.Dimensions != nil {
		d := *req.Dimensions
		if d == 512 {
			dims = 512
		} else if d != 1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "dimensions must be 1024 or 512"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Process texts — check cache first, then TEI for misses
	embeddings := make([][]float32, len(texts))
	promptTokens := 0

	for i, text := range texts {
		// Check cache
		emb, found, err := h.cache.Get(ctx, text, dims)
		if err != nil {
			h.logger.Warn("cache error", "error", err)
		}

		if found && emb != nil {
			h.logger.Info("cache hit", "text_len", len(text), "dim", dims)
			embeddings[i] = emb
			promptTokens += len(text) / 4 // rough estimate
			continue
		}

		// Cache miss — call TEI (always returns 1024d)
		h.logger.Info("cache miss, calling TEI", "text_len", len(text))
		embs, err := h.teiClient.Embed(ctx, []string{text})
		if err != nil {
			h.logger.Error("TEI error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "embedding generation failed"})
			return
		}

		fullEmb := embs[0]

		// Slice to requested dimension
		if dims == 512 && len(fullEmb) >= 512 {
			fullEmb = fullEmb[:512]
		}

		embeddings[i] = fullEmb

		// Cache the result with correct dimension
		if err := h.cache.Set(ctx, text, fullEmb, dims); err != nil {
			h.logger.Warn("cache set error", "error", err)
		}

		promptTokens += len(text) / 4
	}

	// Build response
	data := make([]embed.EmbeddingObj, len(embeddings))
	for i, emb := range embeddings {
		data[i] = embed.EmbeddingObj{
			Object:     "embedding",
			Embedding:  emb,
			Index:      i,
			Dimensions: dims,
		}
	}

	response := embed.EmbeddingsResponse{
		Object: "list",
		Data:   data,
		Model:  h.modelID,
		Usage: embed.Usage{
			PromptTokens: promptTokens,
			TotalTokens:  promptTokens,
		},
	}

	h.logger.Info("embeddings generated", "count", len(embeddings))
	c.JSON(http.StatusOK, response)
}