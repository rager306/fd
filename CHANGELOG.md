# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog 2.0.0](https://keepachangelog.com/en/2.0.0/),
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Repository contributor surface**: `CONTRIBUTING.md` (setup, required gates, PR
  expectations), `SECURITY.md` (supported versions, private vulnerability reporting,
  scope), `CODE_OF_CONDUCT.md` (Contributor Covenant 3.0), issue forms and a pull
  request template under `.github/`, `CODEOWNERS`, `.editorconfig`, `.gitattributes`
  and a weekly `dependabot.yml` for Go modules, workflow actions and container images.
- **Tracked `AGENTS.md`**: the rules for AI coding agents in this repository are now
  versioned instead of local-only. `.gitnexusignore` and the `.gitnexus/` ignore entry
  keep the code-intelligence index consistent for anyone who builds it.

### Fixed

- **`fd_cache_evictions_total` counted nothing.** The counter was registered but
  `ObserveCacheEviction` had no production caller, so it stayed at 0 while the L1
  cache dropped entries on every `FD_CACHE_MAX_SIZE` overflow. `LocalCache` now
  reports each capacity sweep to an eviction observer, `Metrics.ObserveCacheEvictions`
  adds the dropped-entry count, and the metric help states that it counts capacity
  evictions rather than TTL expiry.

### Changed

- **CI hardening**: GitHub Actions are pinned by commit SHA with the version in a
  trailing comment (`actions/checkout` v7.0.1, `actions/setup-go` v7.0.0), and a
  `go vet` step runs alongside tests, lint and `govulncheck`.
- **README**: documents cross-request coalescing (`FD_COALESCE_ENABLED`,
  `FD_COALESCE_WINDOW_MS`), the rate-limit knobs (`FD_RATE_LIMIT_IP_RPM`,
  `FD_RATE_LIMIT_USER_RPM`), the cache-control and queue endpoints, the repository
  layout, and the current verified test, lint and vulnerability counts.

- **L2 Redis key fast path (`REDIS_CACHE_KEY_FASTPATH`, default `true`)**:
  `RedisCache.key()` uses string literals for the closed dimensions set
  (1024/512) instead of `strconv.Itoa`, saving one heap allocation per key
  (~4 B). Key format is identical to the generic branch — verified by
  `TestRedisCacheKeyFastPathMatchesGeneric` — so toggling never splits or
  invalidates the cache population. `false` forces the generic branch.
  Wired through `.env`/`.env.example` and `docker-compose.override.yaml`,
  documented in README and the same-host contract.

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
