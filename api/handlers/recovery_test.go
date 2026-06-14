package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"log/slog"
)

func TestRecoveryMiddlewareCatchesPanicWithEnvelope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/embeddings", http.NoBody)

	logger := slog.New(slog.NewTextHandler(testWriter{t}, nil))
	RecoveryMiddleware(logger)(c)

	// Inject a panic by manually invoking c.Next() with a handler that panics.
	// Since we cannot easily inject a panic into the same context, build
	// a fresh engine that mounts RecoveryMiddleware + a panicking route.
	w2 := httptest.NewRecorder()
	r := gin.New()
	r.Use(RecoveryMiddleware(logger))
	r.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})
	req := httptest.NewRequest(http.MethodGet, "/panic", http.NoBody)
	r.ServeHTTP(w2, req)

	if w2.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500; body=%s", w2.Code, w2.Body.String())
	}
	var resp ErrorResponse
	if err := json.Unmarshal(w2.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, w2.Body.String())
	}
	if resp.Error.Code != CodeInternalError {
		t.Errorf("code = %q, want %q", resp.Error.Code, CodeInternalError)
	}
	if resp.Error.Type != TypeInternalError {
		t.Errorf("type = %q, want %q", resp.Error.Type, TypeInternalError)
	}
	if !strings.Contains(resp.Error.Message, "internal server error") {
		t.Errorf("message = %q, want contains 'internal server error'", resp.Error.Message)
	}
}

func TestRecoveryMiddlewarePassesThroughHappyPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	r := gin.New()
	r.Use(RecoveryMiddleware(slog.New(slog.NewTextHandler(testWriter{t}, nil))))
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	req := httptest.NewRequest(http.MethodGet, "/ok", http.NoBody)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "world") {
		t.Errorf("body = %s, want contains 'world'", w.Body.String())
	}
}

// testWriter is an io.Writer that funnels slog output into t.Log so
// panic stack traces show up in test output for debugging.
type testWriter struct{ t *testing.T }

func (w testWriter) Write(p []byte) (int, error) {
	w.t.Log(strings.TrimSpace(string(p)))
	return len(p), nil
}
