package handlers

import (
	"context"
	"net/http"
	"time"

	"fd-api/embed"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type Embedder interface {
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}

type EmbeddingCache interface {
	GetOrLoad(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error)
}

type EmbeddingsHandler struct {
	teiClient Embedder
	cache     EmbeddingCache
	modelID   string
	logger    *slog.Logger
}

func NewEmbeddingsHandler(teiClient Embedder, c EmbeddingCache, modelID string, logger *slog.Logger) *EmbeddingsHandler {
	return &EmbeddingsHandler{
		teiClient: teiClient,
		cache:     c,
		modelID:   modelID,
		logger:    logger,
	}
}

func (h *EmbeddingsHandler) CreateEmbedding(c *gin.Context) {
	var req embed.EmbeddingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{errorKey: err.Error()})
		return
	}

	texts := req.Input
	if len(texts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{errorKey: "input is required"})
		return
	}

	dims := 1024
	if req.Dimensions != nil {
		d := *req.Dimensions
		if d == 512 {
			dims = 512
		} else if d != 1024 {
			c.JSON(http.StatusBadRequest, gin.H{errorKey: "dimensions must be 1024 or 512"})
			return
		}
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	embeddings := make([][]float32, len(texts))
	promptTokens := 0

	for i, text := range texts {
		emb, err := h.cache.GetOrLoad(ctx, text, dims, func(ctx context.Context) ([]float32, error) {
			embs, err := h.teiClient.Embed(ctx, []string{text})
			if err != nil {
				return nil, err
			}
			return embs[0], nil
		})
		if err != nil {
			h.logger.Error("embedding error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{errorKey: "embedding generation failed"})
			return
		}

		fullEmb := emb
		if dims == 512 && len(fullEmb) >= 512 {
			fullEmb = fullEmb[:512]
		}

		embeddings[i] = fullEmb
		promptTokens += len(text) / 4
	}

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

	c.JSON(http.StatusOK, response)
}
