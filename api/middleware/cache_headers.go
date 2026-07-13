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
	cacheControlValue  = "public, max-age=86400"
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
	var buf [66]byte
	buf[0] = '"'
	hex.Encode(buf[1:65], sum[:])
	buf[65] = '"'
	return string(buf[:])
}

func etagMatches(ifNoneMatch, etag string) bool {
	for _, candidate := range strings.Split(ifNoneMatch, ",") {
		trimmed := strings.TrimSpace(candidate)
		if trimmed == etag || `"`+strings.Trim(trimmed, `"`)+`"` == etag {
			return true
		}
	}
	return false
}
