package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNotFoundMiddlewareWritesEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.NoRoute(NotFoundMiddleware())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/missing", http.NoBody))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", w.Code, w.Body.String())
	}
	var body ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Error.Code != CodeNotFound {
		t.Fatalf("code = %q, want %q", body.Error.Code, CodeNotFound)
	}
	if body.Error.Message != "path /missing not found" {
		t.Fatalf("message = %q", body.Error.Message)
	}
}

func TestMethodNotAllowedMiddlewareWritesEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.HandleMethodNotAllowed = true
	r.NoMethod(MethodNotAllowedMiddleware())
	r.POST("/v1/embeddings", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/embeddings", http.NoBody))

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405; body=%s", w.Code, w.Body.String())
	}
	var body ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Error.Code != CodeMethodNotAllowed {
		t.Fatalf("code = %q, want %q", body.Error.Code, CodeMethodNotAllowed)
	}
	if body.Error.Param != "method" {
		t.Fatalf("param = %q, want method", body.Error.Param)
	}
}
