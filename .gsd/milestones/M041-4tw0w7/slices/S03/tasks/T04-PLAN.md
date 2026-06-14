---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: /metrics endpoint с Prometheus counters/histograms

api/observability/metrics.go: использовать prometheus/client_golang. Metrics: fd_requests_total{status=success|error|timeout} counter, fd_request_duration_seconds histogram (le=0.05/0.1/0.5/1.0/+Inf), fd_batch_size histogram (le=1/10/32/+Inf), fd_errors_total{code=...} counter, fd_model_loaded gauge, fd_cache_hits_total{result=hit|miss} counter (используется в S04). GET /metrics handler использует promhttp.Handler() (text/plain). Middleware MetricsMiddleware оборачивает все requests и инкрементит counters/observations.

## Inputs

- None specified.

## Expected Output

- `api/observability/metrics.go`
- `api/observability/metrics_test.go`

## Verification

Unit tests: после серии requests counter и histogram обновляются корректно. /metrics text/plain содержит все требуемые counter/histogram/gauge. T-H-1..T-H-5 (Section 5.5 existence) pass.
