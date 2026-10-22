package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	headerETag         = "ETag"
	headerIfNoneMatch  = "If-None-Match"
	headerCacheControl = "Cache-Control"
	// private, not public: both cacheable paths (/v1/embeddings, /info) sit
	// behind API key auth and carry caller-specific data. `public` would let a
	// shared cache (proxy, CDN) store the response and serve it across callers.
	// The ETag still enables client-side revalidation.
	cacheControlValue = "private, max-age=86400"
)

// CacheHeaders computes response-body ETags for cacheable fd v2 surfaces.
func CacheHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isCacheHeaderPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		original := c.Writer
		buffered := &cacheHeaderWriter{ResponseWriter: original, body: bytes.Buffer{}}
		c.Writer = buffered
		c.Next()
		c.Writer = original

		status := buffered.Status()
		body := buffered.body.Bytes()
		if status != http.StatusOK {
			original.WriteHeader(status)
			_, _ = original.Write(body)
			return
		}

		etag := responseETag(body)
		original.Header().Set(headerETag, etag)
		original.Header().Set(headerCacheControl, cacheControlValue)
		if etagMatches(c.GetHeader(headerIfNoneMatch), etag) {
			original.WriteHeader(http.StatusNotModified)
			return
		}
		original.WriteHeader(status)
		_, _ = original.Write(body)
	}
}

type cacheHeaderWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *cacheHeaderWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *cacheHeaderWriter) WriteString(data string) (int, error) {
	return w.body.WriteString(data)
}

func isCacheHeaderPath(path string) bool {
	return path == "/v1/embeddings" || path == "/info"
}

func responseETag(body []byte) string {
	sum := sha256.Sum256(body)
	return `"` + hex.EncodeToString(sum[:]) + `"`
}

func etagMatches(ifNoneMatch, etag string) bool {
	// ⚡ Bolt: Zero-allocation fast path for ETag matching by avoiding strings.Split.
	start := 0
	for {
		idx := strings.IndexByte(ifNoneMatch[start:], ',')
		var candidate string
		if idx == -1 {
			candidate = ifNoneMatch[start:]
		} else {
			candidate = ifNoneMatch[start : start+idx]
		}

		trimmed := strings.TrimSpace(candidate)
		if trimmed == etag || `"`+strings.Trim(trimmed, `"`)+`"` == etag {
			return true
		}

		if idx == -1 {
			break
		}
		start = start + idx + 1
	}
	return false
}
