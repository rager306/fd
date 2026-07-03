---
estimated_steps: 24
estimated_files: 6
skills_used: []
---

# T02: Per-stage latency + TEI saturation

Добавить новые метрики в metrics.go:
1. `fd_tei_request_duration_seconds` histogram с bucket set для p95/p99 visibility:
   `[]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.2, 0.5, 1.0, 2.5, 5.0, 10.0}` (TEI inference от 5ms cache-hit до 10s worst-case).
2. `fd_cache_lookup_duration_seconds` histogram с short buckets:
   `[]float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05}` (cache lookup is fast: <50ms).
3. `fd_tei_requests_in_flight` gauge — incremented on TEI call start, decr on end (atomic).
4. `fd_tei_errors_total{reason}` counter — reason values:
   - `timeout` (context.DeadlineExceeded или client timeout)
   - `http_error` (resp.StatusCode != 200)
   - `circuit_open` (TEIClient.checkCircuitOpen returned error)
   - `model_mismatch` (TEI returned wrong count)
   - `transport` (network/dial errors not classifiable выше)

Instrument api/embed/tei.go (TEIClient.doEmbedRequest):
- Start timer in doEmbedRequest (после req construction), measure HTTP roundtrip + body parse
- On error path, classify into one of reason values
- Wrap in closure that captures metrics (TEIClient has metrics field optional, can be nil)

Instrument api/handlers/embeddings.go (CreateEmbedding):
- Wrap cache.GetManyIfPresent call with timer
- Wrap h.teiClient.Embed call with timer (or rely on TEI-side timer)

Embed metrics interface design:
- api/embed/embed.go: добавить optional `WithMetricsObserver()` pattern — TEIClient accepts observability.EmbedderObserver (or keep zero-value tolerance with callback func).

Backward compat: TEIClient constructor signature НЕ меняется. New `WithMetricsObserver(...)` fluent setter. Тесты которые не set observer - nil-safe path.

Files: api/observability/metrics.go (new collectors), api/embed/tei.go (timer + classify + record), api/embed/embed.go (interface extension if needed), api/embed/tei_test.go (new tests for timer classification), api/observability/metrics_test.go (new tests for new metrics), api/main.go (wire observer).

Verify: go test ./api/embed/... ./api/observability/... -race.

## Inputs

- `api/observability/metrics.go`
- `api/embed/tei.go`
- `api/embed/embed.go`

## Expected Output

- `api/observability/metrics.go`
- `api/embed/tei.go`
- `api/embed/embed_test.go (if needed)`
- `api/observability/metrics_test.go`

## Verification

go test ./api/embed/... ./api/observability/... -race ; verify all new metrics names appear in TestMetricsHandlerExposesPrometheusText output.

## Observability Impact

Per-stage latency breakdown visible in /metrics. TEI saturation (in-flight + errors by reason) drives Phase 2 decisions.
