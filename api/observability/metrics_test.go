package observability

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fd-api/handlers"
	"fd-api/lifecycle"
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
		// Phase 0 (M052-mmf99p) instrumentation:
		"fd_cache_entries",
		"fd_cache_memory_bytes",
		"fd_tei_request_duration_seconds",
		"fd_tei_requests_in_flight",
		"fd_tei_errors_total",
		"fd_cache_lookup_duration_seconds",
		"fd_tei_batch_fill_ratio",
	} {
		if !strings.Contains(body, metricName) {
			t.Fatalf("metrics output missing %s:\n%s", metricName, body)
		}
	}
}

func TestMetricsHandlerExposesRuntimeCapacityAndCacheGauges(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := NewMetrics()
	state := lifecycle.NewState()
	done := state.TrackRequest()
	defer done()
	metrics.SetRuntimeObservers(state, 10, func() int { return 3 })
	r := gin.New()
	r.GET("/metrics", metrics.Handler())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Code, w.Body.String())
	}
	body := w.Body.String()
	for _, want := range []string{
		"fd_in_flight_requests 1",
		"fd_in_flight_capacity 10",
		`fd_cache_entries{tier_label="l1"} 3`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("metrics output missing %q:\n%s", want, body)
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
		`fd_cache_hits_total{result="hit",tier_label="all"}`,
		`fd_cache_hits_total{result="miss",tier_label="all"}`,
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

func TestMetricsTEIRequestDurationObserved(t *testing.T) {
	metrics := NewMetrics()
	metrics.ObserveTEIRequestDuration(150 * time.Millisecond)
	metrics.ObserveTEIRequestDuration(25 * time.Millisecond)
	r := gin.New()
	r.GET("/metrics", metrics.Handler())
	body := serveMetricsRequest(r, "GET", "/metrics", "").Body.String()
	for _, want := range []string{
		"fd_tei_request_duration_seconds_bucket",
		"fd_tei_request_duration_seconds_sum",
		"fd_tei_request_duration_seconds_count 2",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in TEI duration metrics:\n%s", want, body)
		}
	}
}

func TestMetricsTEIErrorClassification(t *testing.T) {
	metrics := NewMetrics()
	metrics.IncTEIError("timeout")
	metrics.IncTEIError("http_error")
	metrics.IncTEIError("circuit_open")
	metrics.IncTEIError("transport")
	r := gin.New()
	r.GET("/metrics", metrics.Handler())
	body := serveMetricsRequest(r, "GET", "/metrics", "").Body.String()
	for _, want := range []string{
		`fd_tei_errors_total{reason="timeout"} 1`,
		`fd_tei_errors_total{reason="http_error"} 1`,
		`fd_tei_errors_total{reason="circuit_open"} 1`,
		`fd_tei_errors_total{reason="transport"} 1`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in TEI errors:\n%s", want, body)
		}
	}
}

func TestMetricsTEIInFlightGauge(t *testing.T) {
	metrics := NewMetrics()
	metrics.IncTEIRequestsInFlight()
	metrics.IncTEIRequestsInFlight()
	metrics.DecTEIRequestsInFlight()
	r := gin.New()
	r.GET("/metrics", metrics.Handler())
	body := serveMetricsRequest(r, "GET", "/metrics", "").Body.String()
	if !strings.Contains(body, "fd_tei_requests_in_flight 1") {
		t.Fatalf("TEI in-flight gauge should be 1 after inc/inc/dec:\n%s", body)
	}
}

func TestMetricsBatchFillRatioObserved(t *testing.T) {
	metrics := NewMetrics()
	metrics.ObserveBatchFillRatio(0.25)
	metrics.ObserveBatchFillRatio(0.5)
	metrics.ObserveBatchFillRatio(-1.0) // out of range — no-op
	r := gin.New()
	r.GET("/metrics", metrics.Handler())
	body := serveMetricsRequest(r, "GET", "/metrics", "").Body.String()
	for _, want := range []string{
		"fd_tei_batch_fill_ratio_bucket",
		"fd_tei_batch_fill_ratio_count 2",
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in batch fill:\n%s", want, body)
		}
	}
}

func TestMetricsCacheHitWithTierResult(t *testing.T) {
	metrics := NewMetrics()
	metrics.ObserveCacheResultWithTier("hit", "l1")
	metrics.ObserveCacheResultWithTier("hit", "l2")
	metrics.ObserveCacheResultWithTier("miss", "miss")
	r := gin.New()
	r.GET("/metrics", metrics.Handler())
	body := serveMetricsRequest(r, "GET", "/metrics", "").Body.String()
	for _, want := range []string{
		`fd_cache_hits_total{result="hit",tier_label="l1"} 1`,
		`fd_cache_hits_total{result="hit",tier_label="l2"} 1`,
		`fd_cache_hits_total{result="miss",tier_label="miss"} 1`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in tiered cache hits:\n%s", want, body)
		}
	}
}

func TestMetricsCacheLookupDurationObserved(t *testing.T) {
	metrics := NewMetrics()
	metrics.ObserveCacheLookupDuration(3 * time.Millisecond)
	metrics.ObserveCacheLookupDuration(7 * time.Millisecond)
	r := gin.New()
	r.GET("/metrics", metrics.Handler())
	body := serveMetricsRequest(r, "GET", "/metrics", "").Body.String()
	if !strings.Contains(body, `fd_cache_lookup_duration_seconds_bucket`) {
		t.Fatalf("missing cache lookup duration bucket:\n%s", body)
	}
	if !strings.Contains(body, `fd_cache_lookup_duration_seconds_count 2`) {
		t.Fatalf("cache lookup duration count should be 2:\n%s", body)
	}
}

func TestMetricsRuntimeGaugesIncludeCacheMemory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := NewMetrics()
	state := lifecycle.NewState()
	metrics.SetRuntimeObservers(state, 0, func() int { return 50 })
	r := gin.New()
	r.GET("/metrics", metrics.Handler())
	body := serveMetricsRequest(r, "GET", "/metrics", "").Body.String()
	if !strings.Contains(body, `fd_cache_memory_bytes{tier_label="l1"} 204800`) {
		t.Fatalf("L1 memory: 50 entries × 4096 bytes = 204800:\n%s", body)
	}
}
