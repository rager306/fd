# Continue — fd project handoff

**Last updated:** 2026-07-03 (post v1.0.0 release)
**Branch:** `master` — synced with `origin/master`
**HEAD:** `be67882` — `docs: add CHANGELOG.md (Keep a Changelog 2.0.0), LICENSE (MIT), README v1.0.0 badges`

---

## Current state

fd v1.0.0 released and tagged. All Issue #9 actionable phases complete.
Issue #9 closed on GitHub with summary comment.

### Released milestones (M051–M055)

| Milestone | Title | What |
|---|---|---|
| **M051-h1xr44** | Lifecycle warmup auto-recovery | Background goroutine retries warmup after startup-time TEI unreachability. Env: `FD_WARMUP_START_MAX_ATTEMPTS`, `FD_WARMUP_START_BACKOFF_SEC`, `FD_WARMUP_RECOVERY_INTERVAL_SEC`, `FD_WARMUP_RECOVERY_ENABLED` (default true). |
| **M052-mmf99p** | Phase 0 metrics coverage (Issue #9) | 13 Prometheus series at `/metrics`. Cache by tier, TEI latency/saturation, batch-fill ratio, cache occupancy, lookup duration. Observer patterns in cache + embed packages. |
| **M053-ckvgjw** | Cache capacity tuning (Issue #9 Phase 1a) | `FD_CACHE_MAX_SIZE`, `FD_CACHE_LOCAL_TTL` env vars. `DurationOrDefault` helper. Backward-compatible defaults (10000 entries, 30s TTL). |
| **M054-ybl4jr** | Bulk ingestion queue (Issue #9 Phase 2) | Async `POST /v1/queue` + `GET /v1/queue/:id`. Time-windowed batched worker (10ms × 32 texts). In-memory result store with TTL. Backpressure via bounded channel + 503 + Retry-After. `FD_QUEUE_ENABLED` (default false). 5 queue metrics. |
| **M055-ucu1jl** | Miss-coalescing (Issue #9 Phase 1b) | `CoalescingEmbedder` wraps `embed.Embedder`, merges concurrent calls within `FD_COALESCE_WINDOW_MS` (default 5ms). Capped at 32 texts per batch. `FD_COALESCE_ENABLED` (default false). 44-FZ baseline: 96% reduction in TEI calls under burst. |

### Docker state

| Container | Status | Notes |
|---|---|---|
| `fd_api` | ✅ healthy | Running v1.0.0 binary (rebuilt with `--no-cache`). All M051–M055 features included. |
| `fd_tei` | ✅ healthy | TEI CPU, candle backend, `deepvk/USER-bge-m3` local snapshot. |
| `fd_redis` | ✅ healthy | Namespace `v2`, LRU 2GB. |
| `9router` | ⚠️ unhealthy | Not fd-managed. |
| `spotlight` | Exited | Not fd-managed. |

### Test suite

- **11 Go packages**, all green with `-race`
- New packages: `api/queue/` (worker, store, types), `api/embed/coalescedembedder.go`
- Integration tests: `fd_v2_cache_integration_test.go`, `fd_v2_queue_integration_test.go`
- 44-FZ baseline proof: `api/embed/coalesce_baseline_test.go`

---

## What's done — quick reference

### Env vars added in v1.0.0

| Env var | Default | Purpose |
|---|---|---|
| `FD_WARMUP_START_MAX_ATTEMPTS` | `5` | Startup warmup attempts |
| `FD_WARMUP_START_BACKOFF_SEC` | `5` | Fixed backoff between startup attempts |
| `FD_WARMUP_RECOVERY_INTERVAL_SEC` | `30` | Background recovery retry interval |
| `FD_WARMUP_RECOVERY_ENABLED` | `true` | Feature flag for recovery loop |
| `FD_CACHE_MAX_SIZE` | `10000` | L1 in-memory cache capacity (entries) |
| `FD_CACHE_LOCAL_TTL` | `30s` | L1 entry TTL |
| `FD_QUEUE_ENABLED` | `false` | Feature gate for async `/v1/queue` |
| `FD_QUEUE_MAX_SIZE` | `1024` | Bounded queue capacity |
| `FD_QUEUE_BATCH_MAX_SIZE` | `32` | Max items per worker TEI batch |
| `FD_QUEUE_BATCH_WINDOW_MS` | `10ms` | Worker batch collection window |
| `FD_COALESCE_ENABLED` | `false` | Feature gate for cross-request coalescing |
| `FD_COALESCE_WINDOW_MS` | `5ms` | Coalescing time window |

### New Prometheus metrics

**Cache:** `fd_cache_hits_total{result,tier}`, `fd_cache_entries{tier}`, `fd_cache_memory_bytes{tier}`, `fd_cache_lookup_duration_seconds`

**TEI backend:** `fd_tei_request_duration_seconds`, `fd_tei_requests_in_flight`, `fd_tei_errors_total{reason}`, `fd_tei_batch_fill_ratio`

**Queue:** `fd_queue_depth`, `fd_queue_drain_total`, `fd_queue_submit_total{result}`, `fd_queue_batch_size`, `fd_queue_process_duration_seconds`

### New files (v1.0.0)

```
api/lifecycle/recovery.go + recovery_test.go      (M051)
api/internal/envutil/bool.go + duration.go + tests (M051/M053)
api/observability/metrics.go (extended)            (M052)
api/queue/types.go, id.go, store.go, worker.go     (M054)
api/queue/worker_test.go                           (M054)
api/handlers/queue_handlers.go                     (M054)
api/fd_v2_queue_integration_test.go               (M054)
api/embed/coalescedembedder.go + test              (M055)
api/embed/coalesce_baseline_test.go               (M055)
tools/verify_warmup_recovery.sh                    (M051)
CHANGELOG.md, LICENSE                              (v1.0.0 release)
```

---

## Deferred (from Issue #9)

| Phase | Status | Rationale |
|---|---|---|
| Phase 1c (ORT thread-tuning) | Deferred | TEI uses candle backend, not ORT. Thread config is external. |
| Phase 3 (int8 quantization) | Host-gated | AVX2-only host — regression vs fp16 without VNNI/AMX. |
| Semantic cache (opt-in) | Deferred | KG-exactness risk. Default-off requirement. |
| Phase 2b (Redis-backed queue persistence) | Not started | In-memory queue works; cross-restart durability is future scope. |

---

## Known issues

1. **Docker compose rebuild cache layers** — `docker compose build api` without `--no-cache` may not pick up Go source changes if the Docker layer cache reuses an old `COPY . .` layer. Always use `docker compose build --no-cache api` when verifying new code in the running container.

2. **TEI startup time under CPU contention** — TEI BERT model load on CPU takes 15–20s normally, but up to 3–6 minutes under heavy CPU contention. The lifecycle warmup recovery (M051) handles this automatically; no manual `docker restart fd_api` needed.

3. **env_file directive with `--force-recreate`** — `docker compose up -d --force-recreate api` may drop `FD_API_KEY` from the container env when TEI is not yet healthy at recreate time. Workaround: `docker start fd_api` after TEI becomes healthy, or use `docker compose up -d api` (without `--force-recreate`). Documented in MEM090.

4. **62 closed Jules-bot PRs** — cleaned up on 2026-07-02. All auto-generated Sentinel/Bolt duplicates closed with explanatory comment. Repo is clean.

---

## GSD state

- **Active milestone:** None (all M051–M055 complete)
- **Active slice:** None
- **Queue:** M044-9vahk2 (Upgrade OpenAPI contract to OAS 3.2.0) — planned but not started
- **Issue #9:** Closed on GitHub with summary comment

### Key decisions

- **D052** — Combined warmup approach (startup window + periodic recovery)
- **D053** — Additive tier label for `fd_cache_hits_total` (backward compatible)

### Key memory entries

- **MEM088** — Lifecycle warmup race condition root cause
- **MEM090** — Docker compose env_file force-recreate anomaly

---

## Next steps (candidates)

1. **M044 — Upgrade OpenAPI contract to OAS 3.2.0** — already in QUEUE.md, depends on M041 (✅). Straightforward spec upgrade.

2. **Phase 2b — Redis-backed queue persistence** — extend ResultStore or swap for Redis-backed interface so `/v1/queue` results survive fd-api restarts.

3. **Phase 1c documentation** — add TEI tuning guidance (tokenization-workers, candle thread config) to README Operations section.

4. **GitHub Actions CI** — set up `.github/workflows/` for automated test + build on push/PR, auto-generated release notes for future versions.

5. **Production baseline collection** — deploy v1.0.0, collect real Phase 0 metrics under production traffic, make data-driven decisions for future tuning.

6. **CONTRIBUTING.md + SECURITY.md** — round out open-source community files if the project will accept external contributors.

---

## How to resume

```bash
# Verify everything is in sync
cd /root/fd
git status
git log --oneline -5
docker compose ps

# Run full test suite
cd api && go test ./... -race -count=1

# Smoke the running service
curl http://localhost:8000/health
TOKEN=$(docker exec fd_api printenv FD_API_KEY)
curl -X POST http://localhost:8000/v1/embeddings \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"model":"deepvk/USER-bge-m3","input":"resume smoke"}'
```

### If starting a new milestone

```bash
# GSD workflow
gsd_milestone_generate_id        # get next milestone ID
gsd_plan_milestone               # plan with vision/slices
gsd_plan_slice                   # decompose into tasks
# ... implement, test, verify ...
gsd_task_complete                # close each task
gsd_slice_complete               # close slice
gsd_milestone_validate           # validate
gsd_milestone_complete           # close milestone → auto-commit + push
```

### Key files to read first

- `README.md` — project overview, API, config, operations
- `CHANGELOG.md` — v1.0.0 release notes
- `.gsd/PROJECT.md` — current project state
- `.gsd/STATE.md` — GSD milestone registry
- `api/main.go` — wiring: lifecycle, cache, metrics, queue, coalescer
- `api/lifecycle/recovery.go` — warmup auto-recovery goroutine
- `api/queue/worker.go` — batched queue worker
- `api/embed/coalescedembedder.go` — cross-request coalescer
