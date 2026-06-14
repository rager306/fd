package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"fd-api/buildinfo"
	"fd-api/handlers"

	"github.com/gin-gonic/gin"
)

const headersTestModelID = "deepvk/USER-bge-m3"

func TestHeadersMiddlewareSetsServerAndConnection(t *testing.T) {
	w := serveHeadersRequest(t, "", http.MethodGet, "/ok", "", nil)

	if got := w.Header().Get("Server"); got != "fd/2.0.0" {
		t.Fatalf("Server = %q, want fd/2.0.0", got)
	}
	if got := w.Header().Get("Connection"); got != "keep-alive" {
		t.Fatalf("Connection = %q, want keep-alive", got)
	}
}

func TestHeadersMiddlewareEchoesRequestID(t *testing.T) {
	w := serveHeadersRequest(t, "caller-request-id", http.MethodGet, "/ok", "", nil)

	if got := w.Header().Get(HeaderRequestID); got != "caller-request-id" {
		t.Fatalf("X-Request-Id = %q, want caller-request-id", got)
	}
}

func TestHeadersMiddlewareGeneratesRequestID(t *testing.T) {
	w := serveHeadersRequest(t, "", http.MethodGet, "/ok", "", nil)
	requestID := w.Header().Get(HeaderRequestID)
	if requestID == "" {
		t.Fatal("X-Request-Id should be generated")
	}
	pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !pattern.MatchString(requestID) {
		t.Fatalf("X-Request-Id = %q, want UUIDv4 shape", requestID)
	}
}

func TestHeadersMiddlewareSetsEmbeddingModelAndDimensions(t *testing.T) {
	w := serveHeadersRequest(t, "", http.MethodPost, "/v1/embeddings", `{"model":"test","input":"hello","dimensions":512}`, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	if got := w.Header().Get(HeaderModelID); got != headersTestModelID {
		t.Fatalf("X-Model-Id = %q, want %q", got, headersTestModelID)
	}
	if got := w.Header().Get(HeaderDimensions); got != "512" {
		t.Fatalf("X-Dimensions = %q, want 512", got)
	}
}

func TestHeadersMiddlewareSetsDefaultEmbeddingDimensions(t *testing.T) {
	w := serveHeadersRequest(t, "", http.MethodPost, "/v1/embeddings", `{"model":"test","input":"hello"}`, nil)

	if got := w.Header().Get(HeaderDimensions); got != "1024" {
		t.Fatalf("X-Dimensions = %q, want 1024", got)
	}
}

func TestHeadersMiddlewarePreservesRetryAfterOn503(t *testing.T) {
	w := serveHeadersRequest(t, "", http.MethodGet, "/unready", "", func(c *gin.Context) {
		handlers.WriteErrorWithRetryAfter(c, handlers.CodeModelNotLoaded, "", "warming up", "5")
	})

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusServiceUnavailable, w.Body.String())
	}
	if got := w.Header().Get("Retry-After"); got != "5" {
		t.Fatalf("Retry-After = %q, want 5", got)
	}
	if got := w.Header().Get(HeaderRequestID); got == "" {
		t.Fatal("X-Request-Id should be present on error response")
	}
}

func serveHeadersRequest(t *testing.T, requestID, method, path, body string, customHandler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(HeadersMiddleware(buildinfo.Info{Version: "2.0.0"}, headersTestModelID))
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	if customHandler == nil {
		customHandler = func(c *gin.Context) { c.Status(http.StatusNoContent) }
	}
	r.GET("/unready", customHandler)
	r.POST("/v1/embeddings", ValidateEmbeddingsRequest(), func(c *gin.Context) { c.Status(http.StatusOK) })

	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, http.NoBody)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	if requestID != "" {
		req.Header.Set(HeaderRequestID, requestID)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
