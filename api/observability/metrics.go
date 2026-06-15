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

	runtimeMu        sync.RWMutex
	runtimeState     *lifecycle.State
	runtimeCapacity  int64
	localCacheSizeFn func() int
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
			Help:    "fd HTTP request duration in seconds.",
			Buckets: []float64{0.05, 0.1, 0.5, 1.0},
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
			Help: "Total fd cache lookups by result.",
		}, []string{"result"}),
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
		}, []string{"tier"}),
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
	)
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
	for _, result := range []string{"hit", "miss"} {
		m.cacheHitsTotal.WithLabelValues(result)
	}
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

func (m *Metrics) observeRuntimeGauges() {
	m.runtimeMu.RLock()
	state := m.runtimeState
	capacity := m.runtimeCapacity
	localCacheSizeFn := m.localCacheSizeFn
	m.runtimeMu.RUnlock()
	if state != nil {
		m.inFlightRequests.Set(float64(state.InFlightCount()))
	} else {
		m.inFlightRequests.Set(0)
	}
	m.inFlightCapacity.Set(float64(capacity))
	if localCacheSizeFn != nil {
		m.cacheEntries.WithLabelValues("l1").Set(float64(localCacheSizeFn()))
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
func (m *Metrics) ObserveCacheResult(result string) {
	m.cacheHitsTotal.WithLabelValues(result).Inc()
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
