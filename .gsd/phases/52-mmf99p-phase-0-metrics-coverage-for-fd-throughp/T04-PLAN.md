---
estimated_steps: 16
estimated_files: 3
skills_used: []
---

# T04: Unit tests + edge cases for new metrics

Comprehensive unit test coverage for all new metrics:
1. Existing TestMetricsHandlerExposesPrometheusText: extend list of expected metric names (fd_cache_hits_total{tier=...}, fd_cache_entries{tier=l2}, fd_cache_memory_bytes, fd_tei_request_duration_seconds, fd_cache_lookup_duration_seconds, fd_tei_batch_fill_ratio, fd_tei_requests_in_flight, fd_tei_errors_total).
2. New TestMetricsCacheHitWithTierResult — ObserveCacheResult("l1", "hit") increments correct series, not others.
3. New TestMetricsTEIErrorClassification — observe TEI error with reason="timeout" -> counter increments; "http_error" with status="500" -> another series. Verify label combinations.
4. New TestMetricsTEIRequestDurationObserved — observation populates histogram with custom buckets.
5. New TestMetricsBatchFillRatioObserved — ObserveBatchFillRatio(0.5) populates histogram bucket <=0.5.
6. New TestMetricsTieredCacheLabelsCorrect — integration-style test: setup min TieredCache (with fake local + fake redis), invoke LookUp with miss → hit → miss → l1-hit, check tier labels.
7. New TestMetricsTEIInFlightGauge — start TEI call (use goroutine + sync), gauge increments to 1 during call, decrements after.
8. Edge cases:
   - ObserveCacheResult with unknown tier value — panic? no, log warn.
   - ObserveBatchFillRatio out-of-range (>1.0 or <0.0) — clamp to [0,1] or log warn.
   - TieredCache.SetMetricsObserver(nil) — не падает, no-op.
   - TEIClient.Embed без observer — duration/eerrors не записываются, не падает.

Все тесты используют testing.T, с race detector ( go test -race ).

Files: api/observability/metrics_test.go (extend), api/cache/tiered_test.go (extend or new), api/embed/tei_test.go (extend).

Verify: go test ./api/observability/... ./api/cache/... ./api/embed/... -race -count=1 -v (all pass, no -race warnings). Coverage report: новые metrics paths все покрыты.

## Inputs

- `api/observability/metrics.go`
- `api/cache/tiered.go`
- `api/embed/tei.go`

## Expected Output

- `api/observability/metrics_test.go`
- `api/cache/tiered_test.go`
- `api/embed/tei_test.go`

## Verification

go test ./api/observability/... ./api/cache/... ./api/embed/... -race -count=1 -v (all pass, all new tests included).

## Observability Impact

Tests do not affect runtime observability directly; they verify the contract.
