package handlers

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"math"
	"net/http"
	"time"

	"fd-api/cache"
	"fd-api/embed"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type BatchHandler struct {
	teiClient *embed.TEIClient
	cache     *cache.TieredCache
	modelID   string
	logger    *slog.Logger
}

func NewBatchHandler(teiClient *embed.TEIClient, c *cache.TieredCache, modelID string, logger *slog.Logger) *BatchHandler {
	return &BatchHandler{
		teiClient: teiClient,
		cache:     c,
		modelID:   modelID,
		logger:    logger,
	}
}

func (h *BatchHandler) CreateBatchEmbeddings(c *gin.Context) {
	var req embed.BatchEmbeddingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid batch request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Inputs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "inputs is required"})
		return
	}

	dims := 1024
	if req.Dimensions == 512 {
		dims = 512
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 120*time.Second)
	defer cancel()

	embeddings := make([]string, len(req.Inputs))
	for i, text := range req.Inputs {
		emb, err := h.cache.GetOrLoad(ctx, text, dims, func(ctx context.Context) ([]float32, error) {
			h.logger.Info("batch cache miss, calling TEI", "index", i)
			embs, err := h.teiClient.Embed(ctx, []string{text})
			if err != nil {
				return nil, err
			}
			return embs[0], nil
		})
		if err != nil {
			h.logger.Error("batch embedding error", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "embedding generation failed"})
			return
		}

		fullEmb := emb
		if dims == 512 && len(fullEmb) >= 512 {
			fullEmb = fullEmb[:512]
		}

		embeddings[i] = encodeEmbedding(fullEmb, req.EncodingFormat)
	}

	h.logger.Info("batch embeddings generated", "count", len(req.Inputs))
	c.JSON(http.StatusOK, embed.BatchEmbeddingsResponse{
		Embeddings: embeddings,
		Count:      len(req.Inputs),
		Dimensions: dims,
	})
}

func encodeEmbedding(emb []float32, format string) string {
	if format == "float" {
		b, _ := json.Marshal(emb)
		return string(b)
	}
	return base64.StdEncoding.EncodeToString(float32SliceToBytes(emb))
}

func float32SliceToBytes(slice []float32) []byte {
	b := make([]byte, len(slice)*4)
	for i, v := range slice {
		binary.LittleEndian.PutUint32(b[i*4:], math.Float32bits(v))
	}
	return b
}
