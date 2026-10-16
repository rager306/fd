# Continue — fd project handoff

**Last updated:** 2026-10-16 (post REDIS_CACHE_KEY_FASTPATH)
**Branch:** `master` — synced with `origin/master`
**HEAD:** `8179a42` — `feat(cache): add REDIS_CACHE_KEY_FASTPATH toggle for L2 key fast path`

---

## Current state

fd v1.0.0 released and tagged. All Issue #9 actionable phases complete.
Issue #9 closed on GitHub with summary comment.
CI on master is GREEN (was red 2026-07-03 → 2026-10-15).
No open issues, no open PRs, no remote bot branches (only `master`).

### Released milestones (M051–M055)

| Milestone | Title | What |
|---|---|---|
| **M051-h1xr44** | Lifecycle warmup auto-recovery | Background goroutine retries warmup after startup-time TEI unreachability. Env: `FD_WARMUP_START_MAX_ATTEMPTS`, `FD_WARMUP_START_BACKOFF_SEC`, `FD_WARMUP_RECOVERY_INTERVAL_SEC`, `FD_WARMUP_RECOVERY_ENABLED` (default true). |
| **M052-mmf99p** | Phase 0 metrics coverage (Issue #9) | 13 Prometheus series at `/metrics`. Cache by tier, TEI latency/saturation, batch-fill ratio, cache occupancy, lookup duration. Observer patterns in cache + embed packages. |
| **M053-ckvgjw** | Cache capacity tuning (Issue #9 Phase 1a) | `FD_CACHE_MAX_SIZE`, `FD_CACHE_LOCAL_TTL` env vars. `DurationOrDefault` helper. Backward-compatible defaults (10000 entries, 30s TTL). |
| **M054-ybl4jr** | Bulk ingestion queue (Issue #9 Phase 2) | Async `POST /v1/queue` + `GET /v1/queue/:id`. Time-windowed batched worker (10ms × 32 texts). In-memory result store with TTL. Backpressure via bounded channel + 503 + Retry-After. `FD_QUEUE_ENABLED` (default false). 5 queue metrics. |
| **M055-ucu1jl** | Miss-coalescing (Issue #9 Phase 1b) | `CoalescingEmbedder` wraps `embed.Embedder`, merges concurrent calls within `FD_COALESCE_WINDOW_MS` (default 5ms). Capped at 32 texts per batch. `FD_COALESCE_ENABLED` (default false). 44-FZ baseline: 96% reduction in TEI calls under burst. |

### Post-v1.0.0 commits (2026-08-09 → 2026-10-16)

| Commit | What |
|---|---|
| `bf7b1bc` | fix(api): honor bge m3 token capacity |
| `043467b` | fix(api): clear lint and vet failures that kept CI red (34 golangci findings, `main()` → `run() error` + thin `main`) |
| `7b42906` | fix(api): compare API keys in constant time over SHA-256 digests (length oracle closed) |
| `8ebe847` | fix(queue): do not ignore `crypto/rand` error when minting request ids |
| `4a02331` | chore(deps): Go 1.25 → 1.26, x/net v0.60.0, quic-go v0.59.1 (5 reachable vulns closed, `govulncheck` clean) |
| `8ae6cf2` | perf(cache): `shortHash` via stack buffer (176→64 B, 3→2 allocs) + `hash_test.go` pins 12-hex format |
| `8179a42` | feat(cache): `REDIS_CACHE_KEY_FASTPATH` toggle (default `true`) for L2 key fast path |

Verified 2026-10-16: `go test ./...` 11 pkgs green, `-race` on cache green,
`go vet` clean, `golangci-lint v2.12.2` 0 issues, `govulncheck` 0 vulns,
GitNexus `detect-changes` risk low, Go Quality CI on push success (41s).

### REDIS_CACHE_KEY_FASTPATH (2026-10-16, HEAD)

`RedisCache.key()` uses string literals for the closed dims set (1024/512,
enforced by validation) instead of `strconv.Itoa`: 4→3 allocs, 228→224 B
per key with real `HashText`. Key format is byte-identical to the generic
branch (`TestRedisCacheKeyFastPathMatchesGeneric` pins on/off equivalence),
so toggling never splits the cache population. Strict `ParseBool` in
`RedisCacheOptionsFromEnv` — misconfiguration fails startup. Wired through
`.env`/`.env.example` and `docker-compose.override.yaml` (default `true`).
Local `.env` (gitignored) has it set to `true`.

### Bot-PR hygiene (recurring)

Bolt/Sentinel bots re-create duplicate PRs after each cleanup (5th–6th
copies of the Redis-key fast path arrived 2026-10-15 as #471/#472, closed
2026-10-16 as duplicates with measurements). After closing PRs, also
delete their remote branches — otherwise `bolt-*` branches accumulate on
origin. Backup of old bot branches: `/root/fd-bot-branches.bundle`.

### Docker state

| Container | Status | Notes |
|---|---|---|
| `fd_api` | ✅ healthy | Up 10 days. Image built 2026-08-08 — predates Oct commits (CI fix, auth, Go 1.26, shortHash, KEY_FASTPATH). Rebuild with `--no-cache` before verifying new code live. |
| `fd_tei` | ✅ healthy | TEI CPU, candle backend, `deepvk/USER-bge-m3` local snapshot. |
| `fd_redis` | ✅ healthy | Namespace `v2`, LRU 2GB. |

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

### Env vars added post-v1.0.0

| Env var | Default | Purpose |
|---|---|---|
| `REDIS_CACHE_KEY_FASTPATH` | `true` | L2 key fast path: string literals for dims 1024/512 instead of `strconv.Itoa` (−1 alloc/key). Key format unchanged. `false` forces the generic branch. |

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

4. **Bot-PR cleanup is recurring, not one-off** — 62 Jules-bot PRs closed 2026-07-02,
~400 by 2026-10-15 (backup `/root/fd-bot-branches.bundle`, 403 branches deleted),
then #471/#472 re-arrived the same evening (Redis-key fast-path dupes, closed
2026-10-16). After closing PRs always delete their remote branches too.

---

## GSD state

- **Active milestone:** M044-9vahk2 (Upgrade OpenAPI contract to OAS 3.2.0)
- **Active slice:** S01 (OAS 3.2 delta and validation strategy) — phase `planning`, no DB tasks yet
- **Queue:** M044-9vahk2 only
- **Issues/PRs:** 0 open issues, 0 open PRs (verified 2026-10-16)
- **Current spec:** `api/openapi/spec.go` emits `openapi: 3.1.0`, `spec_test.go` expects `3.1.0`

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
- `CHANGELOG.md` — v1.0.0 release notes + `[Unreleased]` (KEY_FASTPATH entry)
- `.gsd/PROJECT.md` — current project state
- `.gsd/STATE.md` — GSD milestone registry
- `api/main.go` — wiring: lifecycle, cache, metrics, queue, coalescer
- `api/lifecycle/recovery.go` — warmup auto-recovery goroutine
- `api/queue/worker.go` — batched queue worker
- `api/embed/coalescedembedder.go` — cross-request coalescer
