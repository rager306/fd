package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCacheHeadersSetsETagAndCacheControl(t *testing.T) {
	r := cacheHeadersTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get(headerETag); got == "" {
		t.Fatalf("ETag header is empty")
	}
	if got := w.Header().Get(headerCacheControl); got != cacheControlValue {
		t.Fatalf("Cache-Control = %q, want %q", got, cacheControlValue)
	}
	if body := w.Body.String(); body != `{"ok":true}` {
		t.Fatalf("body = %q, want JSON payload", body)
	}
}

func TestCacheHeadersReturnsNotModifiedForMatchingETag(t *testing.T) {
	r := cacheHeadersTestRouter()
	first := httptest.NewRecorder()
	firstReq := httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
	r.ServeHTTP(first, firstReq)
	etag := first.Header().Get(headerETag)
	if etag == "" {
		t.Fatal("first response ETag is empty")
	}

	second := httptest.NewRecorder()
	secondReq := httptest.NewRequest(http.MethodGet, "/info", http.NoBody)
	secondReq.Header.Set(headerIfNoneMatch, etag)
	r.ServeHTTP(second, secondReq)

	if second.Code != http.StatusNotModified {
		t.Fatalf("status = %d, want 304; body=%s", second.Code, second.Body.String())
	}
	if body := second.Body.String(); body != "" {
		t.Fatalf("304 body = %q, want empty", body)
	}
	if got := second.Header().Get(headerETag); got != etag {
		t.Fatalf("ETag = %q, want %q", got, etag)
	}
}

func TestCacheHeadersSkipsNonCacheablePath(t *testing.T) {
	r := cacheHeadersTestRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/live", http.NoBody)
	r.ServeHTTP(w, req)

	if got := w.Header().Get(headerETag); got != "" {
		t.Fatalf("ETag = %q, want empty", got)
	}
}

func cacheHeadersTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CacheHeaders())
	r.GET("/info", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.POST("/v1/embeddings", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	r.GET("/live", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	return r
}
