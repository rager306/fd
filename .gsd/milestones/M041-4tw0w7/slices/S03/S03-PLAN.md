# S03: Observability surface endpoints headers and deep health

**Goal:** Observability surface: endpoints /version, /info, /metrics (Prometheus), /v1/healthcheck, deep /health + обязательные response headers (Server, X-Request-Id, X-Model-Id, X-Dimensions, X-Cache, Retry-After, Connection). Закрывает R-P0-7..R-P0-10, R-P0-11..R-P0-17, R-P1-1, R-P1-2, R-P1-3.
**Demo:** After this, /version, /info, /metrics, /v1/healthcheck, deep /health, /warmup, and response headers are implemented and tested; new exported observability APIs have godoc and pass M043 lint/test/govulncheck gates.

## Must-Haves

- T-H-7..T-H-10 pass: /health deep, /live, /ready, /version
- T-E-1..T-E-3 (Section 5.5 existence) pass: /version 200, /info 200, /metrics 200 text/plain
- T-E-4 pass: /v1/healthcheck 200 (alias)
- T-HDR-1 pass: Server: fd/2.0.0
- T-HDR-2/3 pass: X-Request-Id echo caller-passed OR generated UUIDv4
- T-HDR-4/5 pass: X-Model-Id + X-Dimensions на /v1/embeddings
- T-HDR-8 pass: Retry-After на 503
- T-HDR-9 pass: Connection: keep-alive на всех responses
- /metrics содержит: fd_requests_total{status}, fd_request_duration_seconds histogram, fd_batch_size histogram, fd_errors_total{code}, fd_model_loaded gauge
- Deep /health возвращает status:ok/degraded/down, 200/503 соответственно
- GET /warmup и POST /warmup работают
- golangci-lint pass

## Proof Level

- This slice proves: contract + integration + operational

## Integration Closure

Metrics middleware инкрементит counters/histograms на каждом request. Headers middleware оборачивает все responses. /health deep использует lifecycle state из S02.

## Verification

- Центральный observability slice: добавляет все surfaces (endpoints + headers + metrics) чтобы S04 мог измерять perf и S05 мог репортить features.

## Tasks

- [x] **T01: Added buildinfo metadata package and ldflags wiring for version, build hash, and build date.** `est:2h`
  api/buildinfo package: type Info { Service, Version, Model, ModelVersion, BuildHash, BuildDate, StartedAt, Uptime() time.Duration }. Значения передаются через ldflags при сборке (-X main.Version=2.0.0 -X main.BuildHash=$(git rev-parse --short HEAD) -X main.BuildDate=2026-06-13T00:00:00Z). Default values если ldflags не заданы. Обновить Dockerfile (если нужно) и Makefile (если есть) для передачи ldflags.
  - Files: `api/buildinfo/info.go`, `api/buildinfo/info_test.go`, `Dockerfile`
  - Verify: go test ./api/buildinfo/...: Uptime корректно увеличивается. Build с -ldflags передаёт значения в бинарь.

- [x] **T02: Added /version, /info, and /v1/healthcheck observability endpoints backed by buildinfo and lifecycle state.** `est:2h`
  api/handlers/observability.go: GET /version — возвращает buildinfo.Info + uptime. GET /info — возвращает список моделей с dims=[512,1024], max_input_length_tokens=512, max_batch_size=32, loaded, warmup_done, device (cuda:0/cpu). GET /v1/healthcheck — alias для /health, тот же response. Все endpoints используют lifecycle state из S02.
  - Files: `api/handlers/observability.go`, `api/handlers/observability_test.go`
  - Verify: Integration tests: T-H-10 (/version 200 с version field), T-H-7 (/health deep с model_loaded, warmup_done), T-E-1..T-E-3 (Section 5.5 existence: /version 200, /info 200, /metrics 200 text/plain).

- [x] **T03: Replaced /health with deep lifecycle health reporting ok/degraded/down plus inference and in-flight state.** `est:2h`
  api/handlers/health.go (replace existing): GET /health возвращает { status: ok|degraded|down, time, model_loaded, warmup_done, device, last_inference_at, in_flight_requests }. 200 если status=ok, 503 если degraded/down. last_inference_at обновляется при каждом успешном /v1/embeddings.
  - Files: `api/handlers/health.go`, `api/handlers/health_test.go`
  - Verify: Unit tests: status=ok 200, status=degraded 503, status=down 503. last_inference_at обновляется после inference. T-H-7 pass.

- [x] **T04: Added Prometheus /metrics endpoint and request metrics middleware with counters, histograms, gauges, and error-code labels.** `est:3h`
  api/observability/metrics.go: использовать prometheus/client_golang. Metrics: fd_requests_total{status=success|error|timeout} counter, fd_request_duration_seconds histogram (le=0.05/0.1/0.5/1.0/+Inf), fd_batch_size histogram (le=1/10/32/+Inf), fd_errors_total{code=...} counter, fd_model_loaded gauge, fd_cache_hits_total{result=hit|miss} counter (используется в S04). GET /metrics handler использует promhttp.Handler() (text/plain). Middleware MetricsMiddleware оборачивает все requests и инкрементит counters/observations.
  - Files: `api/observability/metrics.go`, `api/observability/metrics_test.go`
  - Verify: Unit tests: после серии requests counter и histogram обновляются корректно. /metrics text/plain содержит все требуемые counter/histogram/gauge. T-H-1..T-H-5 (Section 5.5 existence) pass.

- [x] **T05: Added response headers middleware for request IDs, server/version, connection, and embedding model/dimensions headers.** `est:3h`
  api/middleware/headers.go: gin middleware. Server: fd/<version> (из buildinfo). X-Request-Id: echo caller-passed X-Request-Id header (любой case), иначе generate UUIDv4. На error pathе — обязательно сохранить X-Request-Id (recovery middleware из S01 должен его прочитать). Connection: keep-alive (default для HTTP/1.1, но explicit). X-Model-Id и X-Dimensions — выставляются на /v1/embeddings responses (model=<deepvk/USER-bge-m3>, dimensions=actual). Retry-After на 429/503 responses.
  - Files: `api/middleware/headers.go`, `api/middleware/headers_test.go`
  - Verify: Unit tests: T-HDR-1 (Server: fd/2.0.0), T-HDR-2/3 (X-Request-Id echo vs generated), T-HDR-4/5 (X-Model-Id + X-Dimensions на /v1/embeddings), T-HDR-8 (Retry-After на 503), T-HDR-9 (Connection: keep-alive).

- [x] **T06: Added /warmup status and trigger endpoints with background pre-warm execution.** `est:2h`
  api/handlers/warmup.go: GET /warmup — { status: ready|warming_up, progress: 0..1 }. POST /warmup — если ready, 200 { status: ready, message: already warm }; если нет, 202 { status: warming_up, message: warmup started } и trigger background warmup (sync.WaitGroup не блокирует).
  - Files: `api/handlers/warmup.go`, `api/handlers/warmup_test.go`
  - Verify: Unit tests: GET /warmup → 200 status:ready после warmup, GET → 200 status:warming_up, progress:<fraction> во время. POST /warmup → 200 если ready, 202 если warming.

- [ ] **T07: Integration tests для endpoints и headers** `est:3h`
  tests/integration/fd_v2_observability_test.go: автоматизировать Section 5.1 T-H-7..T-H-10, Section 5.3 T-HDR-1..T-HDR-10 (кроме T-HDR-6/7 которые зависят от cache в S04), Section 5.5 T-E-1..T-E-3 (endpoints existence). Спека: docs/fd-v2.md Section 5.1, 5.3, 5.5.
  - Files: `tests/integration/fd_v2_observability_test.go`
  - Verify: go test ./tests/integration/... -run TestFdV2Observability -v: все 22 test cases pass.

## Files Likely Touched

- api/buildinfo/info.go
- api/buildinfo/info_test.go
- Dockerfile
- api/handlers/observability.go
- api/handlers/observability_test.go
- api/handlers/health.go
- api/handlers/health_test.go
- api/observability/metrics.go
- api/observability/metrics_test.go
- api/middleware/headers.go
- api/middleware/headers_test.go
- api/handlers/warmup.go
- api/handlers/warmup_test.go
- tests/integration/fd_v2_observability_test.go
