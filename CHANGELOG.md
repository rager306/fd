# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog 2.0.0](https://keepachangelog.com/en/2.0.0/),
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-07-03

First tagged release. fd is a same-host Go embedding API service for
Russian/legal-domain workloads with OpenAI-compatible endpoints, two-tier
caching, and lifecycle observability.

### Added

- **Lifecycle warmup auto-recovery (M051)**: background goroutine retries
  model warmup after startup-time TEI unreachability. Env-configurable via
  `FD_WARMUP_START_MAX_ATTEMPTS`, `FD_WARMUP_START_BACKOFF_SEC`,
  `FD_WARMUP_RECOVERY_INTERVAL_SEC`, `FD_WARMUP_RECOVERY_ENABLED`.
  Eliminates manual `docker restart fd_api` after compose-restart races.

- **Phase 0 metrics coverage (M052, Issue #9)**: 13 new Prometheus series
  at `/metrics` — cache hit-rate by tier (`fd_cache_hits_total{result,tier}`),
  TEI per-stage latency (`fd_tei_request_duration_seconds`), TEI saturation
  (`fd_tei_requests_in_flight`, `fd_tei_errors_total{reason}`), batch-fill
  ratio, cache occupancy (`fd_cache_entries{tier}`, `fd_cache_memory_bytes`),
  cache lookup duration. Observer callback patterns in cache and embed
  packages.

- **Env-tunable L1 cache (M053, Issue #9 Phase 1a)**: `FD_CACHE_MAX_SIZE`
  and `FD_CACHE_LOCAL_TTL` env vars replace hardcoded values.
  `envutil.DurationOrDefault` helper. Backward-compatible defaults.

- **Bulk ingestion queue (M054, Issue #9 Phase 2)**: async `POST /v1/queue`
  + `GET /v1/queue/:id` polling for burst workloads. Time-windowed batched
  worker (10ms × 32 inputs). In-memory result store with TTL eviction.
  Backpressure via bounded channel + 503 + `Retry-After`. Feature gated by
  `FD_QUEUE_ENABLED` (default `false`). 5 new queue metrics.

- **Cross-request embedding coalescer (M055, Issue #9 Phase 1b)**:
  `CoalescingEmbedder` wraps `embed.Embedder` and merges concurrent
  `/v1/embeddings` calls within `FD_COALESCE_WINDOW_MS` (default 5ms).
  Capped at 32 texts per batch to respect TEI `max_client_batch_size`.
  Feature gated by `FD_COALESCE_ENABLED` (default `false`). Synthetic
  baseline on 44-FZ corpus proves 96% reduction in TEI calls under burst.

- **`envutil.DurationOrDefault` and `envutil.BoolOrDefault` helpers** in
  `api/internal/envutil` — reusable env-parsing primitives with safe
  fallbacks.

### Fixed

- `fd_api` no longer stays stuck in `degraded` after a `docker compose
  restart` race where fd starts before TEI finishes loading the BERT model
  on CPU. The lifecycle warmup recovery loop (M051) brings `/health` back
  to `ok` automatically once TEI becomes reachable.

### Changed

- **Request duration histogram buckets** (`fd_request_duration_seconds`)
  expanded from 4 to 7 buckets (5ms–5s) to cover both cache-hot and
  cold-miss latency ranges. Additive — existing dashboards keep working.

- **Cache hits counter** (`fd_cache_hits_total`) gained a new `tier` label
  (`l1`/`l2`/`miss`/`all`). Additive — existing dashboards with
  `by(result)` keep working.

- **L1 cache construction** in `main.go` now reads `FD_CACHE_MAX_SIZE`
  and `FD_CACHE_LOCAL_TTL` instead of hardcoded `10000` entries / `30s`
  TTL. Defaults unchanged — production deploys without env override see
  identical behavior.

### Security

- No security vulnerabilities were found in this release scope.
  Auth middleware remains fail-closed when `FD_API_KEY` is unset.
