package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultCacheDeleteDimensions = 1024

// CacheInvalidator is the cache action surface exposed to authenticated
// operator routes.
type CacheInvalidator interface {
	Delete(ctx context.Context, input string, dim int) error
	Flush(ctx context.Context) (int64, error)
}

// CacheHandler serves authenticated cache invalidation routes.
type CacheHandler struct {
	cache CacheInvalidator
}

// NewCacheHandler returns an authenticated cache invalidation handler.
func NewCacheHandler(cache CacheInvalidator) *CacheHandler {
	return &CacheHandler{cache: cache}
}

type cacheDeleteRequest struct {
	Input      json.RawMessage `json:"input"`
	Dimensions int             `json:"dimensions,omitempty"`
}

type cacheDeleteResponse struct {
	Deleted int `json:"deleted"`
}

type cacheFlushResponse struct {
	Flushed bool  `json:"flushed"`
	Deleted int64 `json:"deleted"`
}

// Flush removes all fd embedding cache entries in the configured namespace.
func (h *CacheHandler) Flush(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	deleted, err := h.cache.Flush(ctx)
	if err != nil {
		WriteError(c, CodeInternalError, "cache", "cache flush failed")
		return
	}
	c.JSON(http.StatusOK, cacheFlushResponse{Flushed: true, Deleted: deleted})
}

// Delete removes cache entries for one or more input texts and a dimension.
func (h *CacheHandler) Delete(c *gin.Context) {
	var req cacheDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteError(c, CodeInvalidJSON, "", "invalid JSON: "+err.Error())
		return
	}
	inputs, ok := parseCacheDeleteInput(req.Input)
	if !ok || len(inputs) == 0 {
		WriteError(c, CodeInputRequired, "input", "input must be a string or array of strings")
		return
	}
	dim := req.Dimensions
	if dim == 0 {
		dim = defaultCacheDeleteDimensions
	}
	if dim != 512 && dim != 1024 {
		WriteError(c, CodeDimensionsInvalid, "dimensions", fmt.Sprintf("dimensions must be 512 or 1024, got %d", dim))
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()
	for _, input := range inputs {
		if err := h.cache.Delete(ctx, input, dim); err != nil {
			WriteError(c, CodeInternalError, "cache", "cache delete failed")
			return
		}
	}
	c.JSON(http.StatusOK, cacheDeleteResponse{Deleted: len(inputs)})
}

func parseCacheDeleteInput(raw json.RawMessage) ([]string, bool) {
	if len(raw) == 0 {
		return nil, false
	}
	var single string
	if err := json.Unmarshal(raw, &single); err == nil {
		single = strings.TrimSpace(single)
		if single == "" {
			return nil, false
		}
		return []string{single}, true
	}
	var many []string
	if err := json.Unmarshal(raw, &many); err != nil {
		return nil, false
	}
	for _, item := range many {
		if strings.TrimSpace(item) == "" {
			return nil, false
		}
	}
	return many, true
}
