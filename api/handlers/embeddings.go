package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"fd-api/embed"
	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
	"log/slog"
)

const (
	// HeaderCache reports whether /v1/embeddings used cache HIT or MISS.
	HeaderCache = "X-Cache"
	cacheHit    = "HIT"
	cacheMiss   = "MISS"
)

// Embedder is the minimal inference interface shared by TEI and ONNX backends.
type Embedder interface {
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}

// EmbeddingCache is the cache surface used by the embeddings handler.
// GetIfPresent is used to peek without triggering a model load, so a
// fully-cached batch can skip the TEI call entirely. Set backfills the
// cache after a model call so future requests for the same text can
// use GetIfPresent. GetOrLoad remains for tests/back-compat.
type EmbeddingCache interface {
	GetIfPresent(ctx context.Context, key string, dim int) ([]float32, bool)
	GetManyIfPresent(ctx context.Context, keys []string, dim int) map[int][]float32
	Set(ctx context.Context, key string, dim int, emb []float32)
	GetOrLoad(ctx context.Context, key string, dim int, loader func(context.Context) ([]float32, error)) ([]float32, error)
}

// EmbeddingsHandler serves the OpenAI-compatible /v1/embeddings endpoint.
type EmbeddingsHandler struct {
	teiClient Embedder
	cache     EmbeddingCache
	modelID   string
	logger    *slog.Logger
}

// NewEmbeddingsHandler wires the embedder, cache, model ID, and logger used by
// the OpenAI-compatible /v1/embeddings endpoint.
func NewEmbeddingsHandler(teiClient Embedder, c EmbeddingCache, modelID string, logger *slog.Logger) *EmbeddingsHandler {
	return &EmbeddingsHandler{
		teiClient: teiClient,
		cache:     c,
		modelID:   modelID,
		logger:    logger,
	}
}

// CreateEmbedding serves POST /v1/embeddings.
//
// Validation (body size, JSON shape, input length, batch size, dimensions,
// encoding_format) is performed by middleware.ValidateEmbeddingsRequest
// and the parsed request is fetched from the gin context. This handler
// therefore only handles embedding generation and response shaping.
//
// All error envelopes go through handlers.WriteError (R-P0-18) so the
// wire shape is OpenAI-compatible and machine-readable.
func (h *EmbeddingsHandler) CreateEmbedding(c *gin.Context) {
	req, ok := h.requestFromContextOrBody(c)
	if !ok {
		return
	}

	dims, encodingFormat := embeddingDefaults(req)
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	embeddings, promptTokens, cacheStatus, ok := h.embeddingsFromCacheOrModel(ctx, req.Input, dims, c)
	if !ok {
		return
	}
	c.Header(HeaderCache, cacheStatus)

	if state, ok := lifecycle.FromContext(c.Request.Context()); ok {
		state.MarkInferenceSuccess()
	}
	c.JSON(http.StatusOK, buildEmbeddingsResponse(embeddings, dims, encodingFormat, h.modelID, promptTokens))
}

func (h *EmbeddingsHandler) requestFromContextOrBody(c *gin.Context) (*embed.EmbeddingsRequest, bool) {
	// Prefer the validated request from middleware; fall back to inline
	// binding so the handler can be mounted standalone in tests or
	// alt entry points without the validation chain. Inline validation
	// is intentionally minimal — production paths must wire
	// middleware.ValidateEmbeddingsRequest.
	if v, ok := c.Get(ContextKeyValidatedRequest); ok {
		if req, ok := v.(*embed.EmbeddingsRequest); ok {
			return req, true
		}
	}

	var inline embed.EmbeddingsRequest
	if err := c.ShouldBindJSON(&inline); err != nil {
		WriteError(c, CodeInvalidJSON, "", "invalid JSON: "+err.Error())
		return nil, false
	}
	if !validateInlineEmbeddingRequest(c, &inline) {
		return nil, false
	}
	return &inline, true
}

func validateInlineEmbeddingRequest(c *gin.Context, req *embed.EmbeddingsRequest) bool {
	if len(req.Input) == 0 {
		WriteError(c, CodeInputRequired, "input", "input is required (non-empty array of strings)")
		return false
	}
	if len(req.Input) > 128 {
		WriteError(c, CodeBatchTooLarge, "input",
			"batch size "+strconv.Itoa(len(req.Input))+" exceeds max 128; split into smaller batches")
		return false
	}
	if req.Dimensions != nil {
		d := *req.Dimensions
		if d != 512 && d != 1024 {
			WriteError(c, CodeDimensionsInvalid, "dimensions",
				"dimensions must be 1024 or 512, got "+strconv.Itoa(d))
			return false
		}
	}
	if req.EncodingFormat != nil {
		ef := *req.EncodingFormat
		if ef != "float" && ef != "base64" {
			WriteError(c, CodeEncodingInvalid, "encoding_format",
				"encoding_format must be float or base64, got \""+ef+"\"")
			return false
		}
	}
	if req.Priority != nil && *req.Priority != "" {
		p := *req.Priority
		if p != "low" && p != "normal" && p != "high" {
			WriteError(c, CodePriorityInvalid, "priority",
				"priority must be low, normal, or high, got \""+p+"\"")
			return false
		}
	}
	return true
}

func embeddingDefaults(req *embed.EmbeddingsRequest) (dims int, encodingFormat string) {
	// Default dimensions: 1024. Pointer absence means "use default".
	dims = 1024
	if req.Dimensions != nil {
		dims = *req.Dimensions
	}

	// Default encoding_format: float. Empty string from nil pointer maps here.
	encodingFormat = embed.EncodingFormatFloat
	if req.EncodingFormat != nil && *req.EncodingFormat != "" {
		encodingFormat = *req.EncodingFormat
	}
	return dims, encodingFormat
}

func (h *EmbeddingsHandler) embeddingsFromCacheOrModel(ctx context.Context, texts []string, dims int, c *gin.Context) (embeddings [][]float32, promptTokens int, cacheStatus string, ok bool) {
	// TEI sub-batch chunking: TEI max_client_batch_size=32 (verified via
	// /info 2026-06-13). When caller sends >32 inputs, split into chunks
	// of 32 and make ONE sequential TEI call per chunk. This avoids the
	// previous work-amplification bug where per-item cache lookups each
	// triggered a full-chunk TEI call.
	const teiSubBatchSize = 32

	embeddings = make([][]float32, len(texts))
	cacheStatus = cacheHit
	for chunkStart := 0; chunkStart < len(texts); chunkStart += teiSubBatchSize {
		chunkEnd := min(chunkStart+teiSubBatchSize, len(texts))
		chunkTokens, chunkCacheMiss, ok := h.fillEmbeddingChunk(ctx, texts, embeddings, chunkStart, chunkEnd, dims, c)
		if !ok {
			return nil, 0, "", false
		}
		if chunkCacheMiss {
			cacheStatus = cacheMiss
		}
		promptTokens += chunkTokens
	}
	return embeddings, promptTokens, cacheStatus, true
}

func (h *EmbeddingsHandler) fillEmbeddingChunk(ctx context.Context, texts []string, embeddings [][]float32, chunkStart, chunkEnd, dims int, c *gin.Context) (promptTokens int, cacheMissed, ok bool) {
	chunk := texts[chunkStart:chunkEnd]
	missIdx, missTexts, promptTokens := h.collectCacheMisses(ctx, chunk, embeddings, chunkStart, dims)
	if len(missTexts) == 0 {
		return promptTokens, false, true
	}

	embs, err := h.teiClient.Embed(ctx, missTexts)
	if err != nil {
		h.logger.Error("embedding error", "error", err,
			"chunk_start", chunkStart, "chunk_end", chunkEnd, "miss_count", len(missTexts))
		WriteError(c, CodeInternalError, "", "embedding generation failed")
		return 0, true, false
	}
	if len(embs) != len(missTexts) {
		h.logger.Error("embedding count mismatch", "expected", len(missTexts), "got", len(embs))
		WriteError(c, CodeInternalError, "", "embedding generation failed: model returned wrong count")
		return 0, true, false
	}

	promptTokens += h.backfillMisses(ctx, missIdx, missTexts, embs, embeddings, chunkStart, dims)
	return promptTokens, true, true
}

func (h *EmbeddingsHandler) collectCacheMisses(ctx context.Context, chunk []string, embeddings [][]float32, chunkStart, dims int) (missIdx []int, missTexts []string, promptTokens int) {
	missIdx = make([]int, 0, len(chunk))
	missTexts = make([]string, 0, len(chunk))
	hits := h.cache.GetManyIfPresent(ctx, chunk, dims)
	for j, text := range chunk {
		if emb, ok := hits[j]; ok {
			embeddings[chunkStart+j] = truncateEmbedding(emb, dims)
			promptTokens += len(text) / 4
			continue
		}
		missIdx = append(missIdx, j)
		missTexts = append(missTexts, text)
	}
	return missIdx, missTexts, promptTokens
}

func (h *EmbeddingsHandler) backfillMisses(ctx context.Context, missIdx []int, missTexts []string, embs, embeddings [][]float32, chunkStart, dims int) (promptTokens int) {
	for k, j := range missIdx {
		emb := embs[k]
		// Store full dim in cache so the next request can short-circuit.
		h.cache.Set(ctx, missTexts[k], dims, emb)
		embeddings[chunkStart+j] = truncateEmbedding(emb, dims)
		promptTokens += len(missTexts[k]) / 4
	}
	return promptTokens
}

func truncateEmbedding(emb []float32, dims int) []float32 {
	if dims == 512 && len(emb) >= 512 {
		return emb[:512]
	}
	return emb
}

func buildEmbeddingsResponse(embeddings [][]float32, dims int, encodingFormat, modelID string, promptTokens int) embed.EmbeddingsResponse {
	data := make([]embed.EmbeddingObj, len(embeddings))
	for i, emb := range embeddings {
		obj := embed.EmbeddingObj{
			Object:     "embedding",
			Index:      i,
			Dimensions: dims,
		}
		if encodingFormat == embed.EncodingFormatBase64 {
			obj.SetBase64(embed.EncodeEmbedding(emb, embed.EncodingFormatBase64))
		} else {
			obj.SetVector(emb)
		}
		data[i] = obj
	}

	return embed.EmbeddingsResponse{
		Object: "list",
		Data:   data,
		Model:  modelID,
		Usage: embed.Usage{
			PromptTokens: promptTokens,
			TotalTokens:  promptTokens,
		},
	}
}
