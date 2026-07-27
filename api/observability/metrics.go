// Package observability exposes Prometheus metrics and gin middleware for fd.
package observability

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"fd-api/embed"
	"fd-api/handlers"
	"fd-api/lifecycle"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	requestStatusSuccess = "success"
	requestStatusError   = "error"
	requestStatusTimeout = "timeout"
)

// LabelTier is the tier label
const LabelTier = "tier"

// Metrics owns fd's Prometheus collectors and registry.
type Metrics struct {
	registry            *prometheus.Registry
	requestsTotal       *prometheus.CounterVec
	requestDuration     prometheus.Histogram
	batchSize           prometheus.Histogram
	errorsTotal         *prometheus.CounterVec
	modelLoaded         prometheus.Gauge
	cacheHitsTotal      *prometheus.CounterVec
	cacheEvictionsTotal prometheus.Counter
	inFlightRequests    prometheus.Gauge
	inFlightCapacity    prometheus.Gauge
	cacheEntries        *prometheus.GaugeVec
	cacheMemoryBytes    *prometheus.GaugeVec

	teiRequestDuration  prometheus.Histogram
	teiRequestsInFlight prometheus.Gauge
	teiErrorsTotal      *prometheus.CounterVec
	cacheLookupDuration prometheus.Histogram
	teiBatchFillRatio   prometheus.Histogram
	queueDepth          prometheus.Gauge
	queueDrainTotal     prometheus.Counter
	queueSubmitTotal    *prometheus.CounterVec
	queueBatchSize      prometheus.Histogram
	queueProcessDuration prometheus.Histogram

	runtimeMu         sync.RWMutex
	runtimeState      *lifecycle.State
	runtimeCapacity   int64
	localCacheSizeFn  func() int
	redisCacheSizeFn  func() int

}

// NewMetrics creates an isolated Prometheus registry with fd collectors.
func NewMetrics() *Metrics {
	metrics := &Metrics{
		registry: prometheus.NewRegistry(),
		requestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fd_requests_total",
			Help: "Total fd HTTP requests by status class.",
		}, []string{"status"}),
		requestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "fd_request_duration_seconds",
			Help:    "fd HTTP request duration in seconds. Covers both cache-hot (1-10ms) and cold-miss (100ms-5s) paths.",
			Buckets: []float64{0.005, 0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
		}),
		batchSize: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "fd_batch_size",
			Help:    "fd embedding request batch size.",
			Buckets: []float64{1, 10, 32},
		}),
		errorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fd_errors_total",
			Help: "Total fd error responses by canonical code.",
		}, []string{"code"}),
		modelLoaded: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "fd_model_loaded",
			Help: "Whether the fd embedding model is loaded and ready (1) or not (0).",
		}),
		cacheHitsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fd_cache_hits_total",
			Help: "Total fd cache lookups by result and tier.",
		}, []string{"result", LabelTier}),
		cacheEvictionsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "fd_cache_evictions_total",
			Help: "Total fd in-memory cache evictions.",
		}),
		inFlightRequests: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "fd_in_flight_requests",
			Help: "Current fd embedding requests in flight.",
		}),
		inFlightCapacity: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "fd_in_flight_capacity",
			Help: "Configured fd embedding in-flight capacity. Zero means unlimited.",
		}),
		cacheEntries: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "fd_cache_entries",
			Help: "Current fd cache entries by tier where cheap to observe.",
		}, []string{LabelTier}),
		cacheMemoryBytes: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "fd_cache_memory_bytes",
			Help: "Approximate memory used by the fd cache by tier. Assumes 1024-dim float32 embeddings (4096 bytes per entry). Not exact — for operational sizing, not billing.",
		}, []string{LabelTier}),

		teiRequestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "fd_tei_request_duration_seconds",
			Help:    "Duration of TEI backend calls in seconds (HTTP round-trip + body parse).",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.5, 1.0, 2.5, 5.0, 10.0},
		}),
		teiRequestsInFlight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "fd_tei_requests_in_flight",
			Help: "Current TEI HTTP requests in flight.",
		}),
		teiErrorsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fd_tei_errors_total",
			Help: "Total TEI error counts by reason.",
		}, []string{"reason"}),
		cacheLookupDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "fd_cache_lookup_duration_seconds",
			Help:    "Time to look up an embedding in the tiered cache (L1 + L2 + backfill).",
			Buckets: []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1},
		}),
		teiBatchFillRatio: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "fd_tei_batch_fill_ratio",
			Help:    "Per-call TEI batch fill ratio in [0.0, 1.0]. Higher means more inputs per TEI HTTP call. Cap denominator is fd's 32-input TEI batch limit.",
			Buckets: []float64{0.0, 0.125, 0.25, 0.375, 0.5, 0.625, 0.75, 0.875, 1.0},
		}),
		queueDepth: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "fd_queue_depth",
			Help: "Current number of items in the async /v1/queue channel.",
		}),
		queueDrainTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "fd_queue_drain_total",
			Help: "Total items drained from the async /v1/queue channel.",
		}),
		queueSubmitTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fd_queue_submit_total",
			Help: "Async /v1/queue submissions by outcome.",
		}, []string{"result"}),
		queueBatchSize: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "fd_queue_batch_size",
			Help:    "Items in each queue-worker TEI batch call.",
			Buckets: []float64{1, 2, 4, 8, 16, 32, 64, 128},
		}),
		queueProcessDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "fd_queue_process_duration_seconds",
			Help:    "Per-batch queue worker processing duration (TEI call + result split).",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1.0, 5.0},
		}),
	}
	metrics.registry.MustRegister(
		metrics.requestsTotal,
		metrics.requestDuration,
		metrics.batchSize,
		metrics.errorsTotal,
		metrics.modelLoaded,
		metrics.cacheHitsTotal,
		metrics.cacheEvictionsTotal,
		metrics.inFlightRequests,
		metrics.inFlightCapacity,
		metrics.cacheEntries,
		metrics.cacheMemoryBytes,
		metrics.teiRequestDuration,
		metrics.teiRequestsInFlight,
		metrics.teiErrorsTotal,
		metrics.cacheLookupDuration,
		metrics.teiBatchFillRatio,
		metrics.queueDepth,
		metrics.queueDrainTotal,
		metrics.queueSubmitTotal,
		metrics.queueBatchSize,
		metrics.queueProcessDuration,
	)
	// init label series for cache_memory_bytes so the gauge appears at
	// startup rather than only after first data point, keeping /metrics
	// shape predictable even on fresh deployments.
	for _, tier := range []string{"l1", "l2"} {
		metrics.cacheMemoryBytes.WithLabelValues(tier)
	}
	metrics.initLabelSeries()
	return metrics
}

func (m *Metrics) initLabelSeries() {
	for _, status := range []string{requestStatusSuccess, requestStatusError, requestStatusTimeout} {
		m.requestsTotal.WithLabelValues(status)
	}
	for _, code := range handlers.AllErrorCodes() {
		m.errorsTotal.WithLabelValues(code)
	}
	// Cache hit-rate by tier: pre-register all result x tier combinations so
	// /metrics exposes the full series even on a cold start. Existing
	// result-only series stay compatible (additive tier label).
	for _, result := range []string{"hit", "miss"} {
		for _, tier := range []string{"l1", "l2", "miss", "all"} {
			m.cacheHitsTotal.WithLabelValues(result, tier)
		}
	}
	for _, reason := range []string{"timeout", "http_error", "circuit_open", "model_mismatch", "transport"} {
		m.teiErrorsTotal.WithLabelValues(reason)
	}
	// Async /v1/queue: submit outcome label cardinality.
	for _, result := range []string{"accepted", "rejected"} {
		m.queueSubmitTotal.WithLabelValues(result)
	}
}

// ObserveTEIRequestDuration records the duration of a single TEI call.
// Called by TEIClient.doEmbedRequest at the end of each call (success or
// failure). Cheap (~10ns) by design.
func (m *Metrics) ObserveTEIRequestDuration(d time.Duration) {
	m.teiRequestDuration.Observe(d.Seconds())
}

// IncTEIRequestsInFlight and DecTEIRequestsInFlight manage a gauge for
// concurrent TEI calls. Safe to call from multiple goroutines.
func (m *Metrics) IncTEIRequestsInFlight() { m.teiRequestsInFlight.Inc() }
// DecTEIRequestsInFlight decrements the in-flight requests count.
func (m *Metrics) DecTEIRequestsInFlight() { m.teiRequestsInFlight.Dec() }

// IncTEIError records a TEI error with a canonical reason label.
func (m *Metrics) IncTEIError(reason string) {
	m.teiErrorsTotal.WithLabelValues(reason).Inc()
}

// ObserveCacheLookupDuration records the duration of a single cache lookup
// (L1 + L2 + backfill). Called by TieredCache.GetManyIfPresent.
func (m *Metrics) ObserveCacheLookupDuration(d time.Duration) {
	m.cacheLookupDuration.Observe(d.Seconds())
}

// ObserveBatchFillRatio records the ratio of inputs in a TEI call relative
// to the 32-input cap. Inputs above the cap are chunked in handlers.obs
func (m *Metrics) ObserveBatchFillRatio(ratio float64) {
	if ratio < 0 || ratio > 1 {
		return // bucket is bounded, skip out-of-range observations
	}
	m.teiBatchFillRatio.Observe(ratio)
}

// Handler returns a Prometheus text-format HTTP handler for /metrics.
func (m *Metrics) Handler() gin.HandlerFunc {
	handler := promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
	return func(c *gin.Context) {
		m.observeRuntimeGauges()
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

// Middleware records request counters, duration, batch size, error codes, and
// lifecycle model-loaded gauge values after downstream handlers complete.
func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()

		statusCode := c.Writer.Status()
		m.requestsTotal.WithLabelValues(requestStatus(statusCode)).Inc()
		m.requestDuration.Observe(time.Since(started).Seconds())
		m.observeBatchSize(c)
		m.observeErrorCode(c, statusCode)
		m.observeModelLoaded(c)
	}
}

// SetRuntimeObservers configures cheap runtime gauges collected at scrape time.
func (m *Metrics) SetRuntimeObservers(state *lifecycle.State, capacity int64, localCacheSizeFn func() int) {
	m.runtimeMu.Lock()
	defer m.runtimeMu.Unlock()
	m.runtimeState = state
	m.runtimeCapacity = capacity
	m.localCacheSizeFn = localCacheSizeFn
}

// SetRedisCacheSizeObserver registers a callback returning the L2 Redis
// namespace size. Called from observeRuntimeGauges at scrape cadence.
// Safe to leave nil; the gauge stays at zero.
func (m *Metrics) SetRedisCacheSizeObserver(fn func() int) {
	m.runtimeMu.Lock()
	defer m.runtimeMu.Unlock()
	m.redisCacheSizeFn = fn
}

// Approximate memory used per embedding cache entry: 4096 bytes for
// 1024-dim float32 (1024 * 4) + a small overhead for hash key + LRU bookkeeping.
const approxEmbeddingBytes = 4096

func (m *Metrics) observeRuntimeGauges() {
	m.runtimeMu.RLock()
	state := m.runtimeState
	capacity := m.runtimeCapacity
	localCacheSizeFn := m.localCacheSizeFn
	redisCacheSizeFn := m.redisCacheSizeFn
	m.runtimeMu.RUnlock()
	if state != nil {
		m.inFlightRequests.Set(float64(state.InFlightCount()))
	} else {
		m.inFlightRequests.Set(0)
	}
	m.inFlightCapacity.Set(float64(capacity))
	if localCacheSizeFn != nil {
		n := localCacheSizeFn()
		m.cacheEntries.WithLabelValues("l1").Set(float64(n))
		m.cacheMemoryBytes.WithLabelValues("l1").Set(float64(n) * approxEmbeddingBytes)
	} else {
		m.cacheEntries.WithLabelValues("l1").Set(0)
		m.cacheMemoryBytes.WithLabelValues("l1").Set(0)
	}
	if redisCacheSizeFn != nil {
		n := redisCacheSizeFn()
		m.cacheEntries.WithLabelValues("l2").Set(float64(n))
		m.cacheMemoryBytes.WithLabelValues("l2").Set(float64(n) * approxEmbeddingBytes)
	}
}

// SetModelLoaded updates fd_model_loaded explicitly for non-request lifecycle changes.
func (m *Metrics) SetModelLoaded(loaded bool) {
	if loaded {
		m.modelLoaded.Set(1)
		return
	}
	m.modelLoaded.Set(0)
}

// ObserveCacheResult increments fd_cache_hits_total for future cache middleware.
//
// The call uses an "all" tier sentinel so legacy single-arg callers
// (including integration tests) keep contributing to a stable series
// without needing to know the new tier label. Prefer
// ObserveCacheResultWithTier in new code.
func (m *Metrics) ObserveCacheResult(result string) {
	m.cacheHitsTotal.WithLabelValues(result, "all").Inc()
}

// ObserveCacheResultWithTier increments fd_cache_hits_total with tier label
// in addition to result. Tier must be "l1", "l2", "miss", or "all".
func (m *Metrics) ObserveCacheResultWithTier(result, tier string) {
	m.cacheHitsTotal.WithLabelValues(result, tier).Inc()
}

// ObserveCacheEviction increments fd_cache_evictions_total.
func (m *Metrics) ObserveCacheEviction() {
	m.cacheEvictionsTotal.Inc()
}

func (m *Metrics) observeBatchSize(c *gin.Context) {
	value, ok := c.Get(handlers.ContextKeyValidatedRequest)
	if !ok {
		return
	}
	req, ok := value.(*embed.EmbeddingsRequest)
	if !ok {
		return
	}
	m.batchSize.Observe(float64(len(req.Input)))
}

func (m *Metrics) observeErrorCode(c *gin.Context, statusCode int) {
	if statusCode < http.StatusBadRequest {
		return
	}
	code, ok := c.Get(handlers.ContextKeyErrorCode)
	if ok {
		if codeValue, ok := code.(string); ok && codeValue != "" {
			m.errorsTotal.WithLabelValues(codeValue).Inc()
			return
		}
	}
	m.errorsTotal.WithLabelValues(strconv.Itoa(statusCode)).Inc()
}

func (m *Metrics) observeModelLoaded(c *gin.Context) {
	state, ok := lifecycle.FromContext(c.Request.Context())
	if !ok {
		return
	}
	m.SetModelLoaded(state.IsReady())
}

func requestStatus(statusCode int) string {
	if statusCode == http.StatusGatewayTimeout {
		return requestStatusTimeout
	}
	if statusCode >= http.StatusBadRequest {
		return requestStatusError
	}
	return requestStatusSuccess
}

// SetQueueDepth sets the queue depth gauge. Called from the queue handler at
// submit time or from a periodic scrape-time observer.
func (m *Metrics) SetQueueDepth(n int) {
	m.queueDepth.Set(float64(n))
}

// IncQueueDrain increments fd_queue_drain_total by n.
func (m *Metrics) IncQueueDrain(n int) {
	m.queueDrainTotal.Add(float64(n))
}

// IncQueueSubmit records a /v1/queue submission outcome. result must be
// "accepted" or "rejected".
func (m *Metrics) IncQueueSubmit(result string) {
	m.queueSubmitTotal.WithLabelValues(result).Inc()
}

// ObserveQueueBatchSize records the size of a worker-processed batch.
func (m *Metrics) ObserveQueueBatchSize(n int) {
	m.queueBatchSize.Observe(float64(n))
}

// ObserveQueueProcessDuration records the duration of a single worker batch
// (TEI call + result split).
func (m *Metrics) ObserveQueueProcessDuration(d time.Duration) {
	m.queueProcessDuration.Observe(d.Seconds())
}
