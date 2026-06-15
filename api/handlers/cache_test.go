package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeCacheInvalidator struct {
	deleted []struct {
		input string
		dim   int
	}
	flushCount int
	flushN     int64
	err        error
}

func (f *fakeCacheInvalidator) Delete(ctx context.Context, input string, dim int) error {
	if f.err != nil {
		return f.err
	}
	f.deleted = append(f.deleted, struct {
		input string
		dim   int
	}{input: input, dim: dim})
	return nil
}

func (f *fakeCacheInvalidator) Flush(ctx context.Context) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	f.flushCount++
	return f.flushN, nil
}

func serveCacheRequest(handler *CacheHandler, path, body string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/cache/flush", handler.Flush)
	r.POST("/v1/cache/delete", handler.Delete)

	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCacheHandlerFlushReportsDeletedCount(t *testing.T) {
	invalidator := &fakeCacheInvalidator{flushN: 7}
	w := serveCacheRequest(NewCacheHandler(invalidator), "/v1/cache/flush", "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if invalidator.flushCount != 1 {
		t.Fatalf("flushCount = %d, want 1", invalidator.flushCount)
	}
	var resp struct {
		Flushed bool  `json:"flushed"`
		Deleted int64 `json:"deleted"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !resp.Flushed || resp.Deleted != 7 {
		t.Fatalf("response = %+v, want flushed=true deleted=7", resp)
	}
}

func TestCacheHandlerDeleteAcceptsStringAndArrayInput(t *testing.T) {
	invalidator := &fakeCacheInvalidator{}
	body := `{"input":["first","second"],"dimensions":512}`
	w := serveCacheRequest(NewCacheHandler(invalidator), "/v1/cache/delete", body)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if len(invalidator.deleted) != 2 {
		t.Fatalf("deleted = %d, want 2", len(invalidator.deleted))
	}
	if invalidator.deleted[0].input != "first" || invalidator.deleted[0].dim != 512 {
		t.Fatalf("first delete = %+v", invalidator.deleted[0])
	}
	if invalidator.deleted[1].input != "second" || invalidator.deleted[1].dim != 512 {
		t.Fatalf("second delete = %+v", invalidator.deleted[1])
	}
}

func TestCacheHandlerDeleteDefaultsDimensionsTo1024(t *testing.T) {
	invalidator := &fakeCacheInvalidator{}
	w := serveCacheRequest(NewCacheHandler(invalidator), "/v1/cache/delete", `{"input":"hello"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if got := invalidator.deleted[0].dim; got != 1024 {
		t.Fatalf("dim = %d, want 1024", got)
	}
}

func TestCacheHandlerDeleteRejectsMalformedInput(t *testing.T) {
	w := serveCacheRequest(NewCacheHandler(&fakeCacheInvalidator{}), "/v1/cache/delete", `{"input":[123]}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body=%s", w.Code, w.Body.String())
	}
}

func TestCacheHandlerMapsInvalidatorErrors(t *testing.T) {
	w := serveCacheRequest(NewCacheHandler(&fakeCacheInvalidator{err: errors.New("boom")}), "/v1/cache/flush", "")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w.Code, w.Body.String())
	}
}
