---
id: T02
parent: S01
milestone: M052-mmf99p
key_files:
  - api/observability/metrics.go
  - api/embed/tei.go
  - api/main.go
key_decisions:
  - Observers struct с nil-tolerant fields (no-op) — TEIClient.safe to share with tests without observers
  - Per-call timing в doEmbedRequest через defer — captures success and error uniformly
  - error classification через strings.Contains на err msg — pragmatic, no need to add sentinel error types
  - fd_tei_requests_in_flight как gauge (Inc/Dec) — дополняет fd_in_flight_requests (lifecycle gate) разделением semantic: lifecycle gate vs TEI HTTP calls в полёте
duration: 
verification_result: passed
completed_at: 2026-07-03T09:07:56.711Z
blocker_discovered: false
---

# T02: TEI per-stage latency + saturation: fd_tei_request_duration_seconds (11 buckets), fd_tei_requests_in_flight gauge, fd_tei_errors_total{reason} + TEIClient.Observers struct + classifyTEIError. Все 10 пакетов тестов зелёные с -race

**TEI per-stage latency + saturation: fd_tei_request_duration_seconds (11 buckets), fd_tei_requests_in_flight gauge, fd_tei_errors_total{reason} + TEIClient.Observers struct + classifyTEIError. Все 10 пакетов тестов зелёные с -race**

## What Happened

T02 closes per-stage TEI observability:

1. **metrics.go** extended with:
   - `fd_tei_request_duration_seconds` histogram, 11 buckets from 5ms to 10s — covers cache-hit TEI (5-20ms) through worst-case cold (3+s)
   - `fd_tei_requests_in_flight` gauge
   - `fd_tei_errors_total{reason}` counter, 5 reason labels: timeout, http_error, circuit_open, model_mismatch, transport
   - `fd_cache_lookup_duration_seconds` histogram (registered but not yet wired in tiered.go - deferred)
   - `fd_tei_batch_fill_ratio` histogram (registered but not yet wired - T03)
   - `ObserveTEIRequestDuration`, `IncTEIRequestsInFlight`, `DecTEIRequestsInFlight`, `IncTEIError`, `ObserveCacheLookupDuration`, `ObserveBatchFillRatio` methods

2. **TEIClient.Observers struct** — optional callback fields:
   - `ObserveDuration func(time.Duration)`
   - `ObserveError func(reason string)`
   - `IncInFlight`, `DecInFlight funcs()`
   - `ObserveBatchFill func(inputs int)`
   - nil fields are no-op (safe for tests)
   - `WithObservers(obs)` fluent setter returns receiver for chaining

3. **TEIClient.doEmbedRequest** instrumented:
   - `started := time.Now()` + defer `recordDuration(started, err)` for both success and error paths
   - Inc/Dec in-flight via defer
   - `classifyTEIError(err)` maps to canonical reason: errors.Is(context.DeadlineExceeded/Canceled) -> timeout; err msg contains 'TEI returned status' -> http_error; 'circuit open' -> circuit_open; 'wrong count' -> model_mismatch; else transport

4. **main.go wiring** — moved to after `metrics := observability.NewMetrics()`:
```go
teiClient.WithObservers(embed.Observers{
    ObserveDuration:  metrics.ObserveTEIRequestDuration,
    ObserveError:     metrics.IncTEIError,
    IncInFlight:      metrics.IncTEIRequestsInFlight,
    DecInFlight:      metrics.DecTEIRequestsInFlight,
    ObserveBatchFill: func(n int) { metrics.ObserveBatchFillRatio(float64(n) / 32.0) },
})
```

5. **All 10 packages green with -race**.

fc - всё работает!

## Verification

go test ./api/... -race -count=1 — все 10 пакетов зелёные. New metrics: fd_tei_request_duration_seconds (11 buckets), fd_tei_requests_in_flight gauge, fd_tei_errors_total{reason} counter (timeout|http_error|circuit_open|model_mismatch|transport). TEIClient обогащён Observers struct + WithObservers method + classifyTEIError + recordDuration timer.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go build ./... && go vet ./...` | 0 | pass | 3200ms |
| 2 | `cd /root/fd/api && go test ./... -race -count=1 -timeout 120s` | 0 | pass | 13000ms |

## Deviations

Wiring in main.go placement — сначала поставил metrics reference до `metrics := observability.NewMetrics()`, переместил wiring после создания metrics object. Wiring is correct now.

## Known Issues

cache_lookup_duration histogram зарегистрирован (T02) но в TieredCache.GetManyIfPresent ещё не вызывается — это оставлено для T03 (batch-fill ratio) или отдельного T02-extension. Сейчас cache_lookup_duration имеет 0 series.

## Files Created/Modified

- `api/observability/metrics.go`
- `api/embed/tei.go`
- `api/main.go`
