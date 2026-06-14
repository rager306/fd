package observability

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fd-api/handlers"
	"fd-api/middleware"

	"github.com/gin-gonic/gin"
)

func TestMetricsHandlerExposesPrometheusText(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := NewMetrics()
	r := gin.New()
	r.GET("/metrics", metrics.Handler())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", w.Code, http.StatusOK, w.Body.String())
	}
	if contentType := w.Header().Get("Content-Type"); !strings.Contains(contentType, "text/plain") {
		t.Fatalf("Content-Type = %q, want text/plain", contentType)
	}
	body := w.Body.String()
	for _, metricName := range []string{
		"fd_requests_total",
		"fd_request_duration_seconds",
		"fd_batch_size",
		"fd_errors_total",
		"fd_model_loaded",
		"fd_cache_hits_total",
		"fd_cache_evictions_total",
	} {
		if !strings.Contains(body, metricName) {
			t.Fatalf("metrics output missing %s:\n%s", metricName, body)
		}
	}
}

func TestMetricsMiddlewareRecordsRequestsErrorsAndBatchSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := NewMetrics()
	r := gin.New()
	r.Use(metrics.Middleware())
	r.GET("/ok", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	r.GET("/fail", func(c *gin.Context) {
		handlers.WriteError(c, handlers.CodeModelNotLoaded, "", "warming up")
	})
	r.POST("/v1/embeddings", middleware.ValidateEmbeddingsRequest(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	r.GET("/metrics", metrics.Handler())

	serveMetricsRequest(r, http.MethodGet, "/ok", "")
	serveMetricsRequest(r, http.MethodGet, "/fail", "")
	serveMetricsRequest(r, http.MethodPost, "/v1/embeddings", `{"model":"test","input":["a","b"]}`)

	w := serveMetricsRequest(r, http.MethodGet, "/metrics", "")
	body := w.Body.String()
	for _, want := range []string{
		`fd_requests_total{status="success"}`,
		`fd_requests_total{status="error"}`,
		`fd_errors_total{code="model_not_loaded"}`,
		"fd_request_duration_seconds_bucket",
		"fd_batch_size_bucket",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics output missing %s:\n%s", want, body)
		}
	}
}

func TestMetricsModelLoadedAndCacheResult(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := NewMetrics()
	metrics.SetModelLoaded(true)
	metrics.ObserveCacheResult("hit")
	metrics.ObserveCacheResult("miss")
	metrics.ObserveCacheEviction()
	r := gin.New()
	r.GET("/metrics", metrics.Handler())

	w := serveMetricsRequest(r, http.MethodGet, "/metrics", "")
	body := w.Body.String()
	for _, want := range []string{
		"fd_model_loaded 1",
		`fd_cache_hits_total{result="hit"}`,
		`fd_cache_hits_total{result="miss"}`,
		"fd_cache_evictions_total 1",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics output missing %s:\n%s", want, body)
		}
	}
}

func serveMetricsRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	var req *http.Request
	if body == "" {
		req = httptest.NewRequest(method, path, http.NoBody)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
