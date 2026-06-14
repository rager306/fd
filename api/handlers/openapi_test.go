package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestOpenAPIHandlerReturnsJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/openapi.json", NewOpenAPIHandler())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", http.NoBody)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	var spec map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &spec); err != nil {
		t.Fatalf("unmarshal spec: %v", err)
	}
	if spec["openapi"] != "3.1.0" {
		t.Fatalf("openapi = %v, want 3.1.0", spec["openapi"])
	}
}

func TestDocsHandlerReturnsSwaggerUIHTML(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/docs", NewDocsHandler())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/docs", http.NoBody)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	if got := w.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", got)
	}
	body := w.Body.String()
	if !strings.Contains(body, "swagger-ui") || !strings.Contains(body, "/openapi.json") {
		t.Fatalf("docs body missing swagger-ui or /openapi.json: %s", body)
	}
}
