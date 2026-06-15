# Issue #8: Agent-native audit: cache-flush surface, health context (last_error + deps), metrics capacity, traces hardening, config-extraction, dead publicMetrics

State: OPEN
URL: https://github.com/rager306/fd/issues/8
Labels: enhancement
Author: rager306
Created: 2026-06-15T12:25:56Z
Updated: 2026-06-15T12:25:56Z

## Body

# Agent-native architecture audit — actionable gaps

Follow-up to #3 / #6 / #7. Ran the 8-principle agent-native audit (`ce-agent-native-audit`, 8 parallel reviewers) against current `master` (HEAD `687d6dd`). Build green, tests green. **Overall ~70% on the 6 applicable principles** (UI Integration and Prompt-Native are NA — fd is a headless deterministic inference API, correctly). Strongest: Shared Workspace (100%), Capability Discovery (93%), Context Injection (75%). This issue tracks the **actionable agent-native gaps** — places where an AI agent/client driving the service can't observe or control something it reasonably should.

> Scope note: the OpenAPI `/embeddings/batch` spec-drift item is **intentionally excluded** from this issue per request.

---

## HIGH impact

### AN-A — No cache flush / invalidation surface (CRUD + Action Parity + the original AN-01)
**Where:** `api/cache/local.go:91` (`LocalCache.Delete` — defined but **zero production callers**), `api/cache/tiered.go` (no `Delete`/`Flush`), `api/cache/redis.go` (no `Del`/`FlushDB`), no `DELETE /v1/cache` or `POST /cache/flush` route in `api/main.go`.
**Why it matters (agent angle):** an agent that knows the model was retrained, or that a cached vector is wrong, has **no way to purge stale vectors** without restarting the process or reaching into Redis out-of-band. The cache-aside backfill (`TieredCache`, `loadBatchEmbeddings`) is shared across all callers, so a single bad entry poisons everyone until TTL. This closes the orphaned-`Delete` finding from the prior review (still orphaned) and is the single highest-leverage agent-native gap.
**Fix:** add `Delete(ctx, key)` + `Flush(ctx)` to `TieredCache` (delegate to L1 `Delete` + a new `RedisCache.Del`), then expose 1–2 admin-guarded routes — `DELETE /v1/cache/:keyHash` and/or `POST /v1/cache/flush`. Gate behind the existing `APIKeyAuthFromEnv`. Wire the existing dead `LocalCache.Delete`.

### AN-B — `/health` hides `last_error` detail and dependency (TEI/Redis) health (Context Injection)
**Where:** `api/lifecycle/state.go` (`LastError()` is tracked and gates readiness, but never serialized); `api/handlers/health.go` (`DeepHealthResponse`); `api/handlers/probes.go:15` (`/live` explicitly skips downstream).
**Why it matters:** an agent sees `status: degraded` / `/ready` 503 with **no WHY**, and cannot tell whether the cause is fd itself, TEI being down, or Redis unreachable. This blocks autonomous recovery — the agent can't route around the right dependency. The state is already held internally; it's just not exposed.
**Fix:** (1) add `last_error {code, message, at}` to `DeepHealthResponse` when status ≠ ok; (2) add a `dependencies {tei:{reachable,latency_ms}, redis:{reachable,latency_ms,namespace}}` block (reuse the existing redis-ping path + a lightweight TEI probe).

---

## MEDIUM impact

### AN-C — `/metrics` lacks capacity ceiling and cache-occupancy (Context Injection)
**Where:** `api/middleware/lifecycle.go` (`maxInFlight` enforced, returns `model_overloaded`, but the ceiling is never reported); `api/handlers/health.go` (`in_flight_requests` shown, limit hidden); `api/observability/metrics.go` (counters for hit/miss/eviction exist, no size gauge).
**Why it matters:** an agent seeing `in_flight_requests: 5` can't judge utilization (5/10 vs 5/10000), and can't reason about cache pressure (no current-entries/size gauge). It can't decide whether to back off.
**Fix:** expose `in_flight_capacity` alongside `in_flight_requests` in `/health`; emit `fd_in_flight_requests`/`fd_cache_entries{l1,l2}`/`fd_cache_size_bytes` gauges to `/metrics`.

### AN-D — `/v1/traces` exposes every caller's requests to any single bearer holder (Shared Workspace)
**Where:** `api/observability/traces.go` (`Snapshot` returns last 100 across all callers, no per-caller filter), `api/main.go:331` (`GET /v1/traces` protected only by the single shared `FD_API_KEY`).
**Why it matters:** fd is single-tenant by design (`fd-v2.md`), but if >1 independent client ever shares an instance, `/v1/traces` leaks their request paths/latencies/volumes to each other. Worth a hardening knob before that happens.
**Fix:** gate `/v1/traces` behind a distinct admin token or `FD_TRACES_PUBLIC=false`; or add per-caller scoping if/when a caller-id exists.

### AN-E — Deployment-tunable policy baked into handlers/middleware (Tools as Primitives)
**Where:** `api/middleware/validation.go` (unexported consts `maxBatchSize=128`, `maxInputChars=2048`, `maxTotalInputChars`, `maxRequestBodyBytes=10MB`, dim allowlist `{512,1024}`); `api/handlers/embeddings.go` (`teiSubBatchSize=32` hardcoded inline); per-handler timeouts (30s/120s) hardcoded.
**Why it matters:** these are the values an operator/deployment would want to tune, but they require source edits + rebuild. The components are otherwise good primitives; this is "configuration as code."
**Fix:** lift bounds into a `ValidationLimits` struct + `ValidateEmbeddingsRequestFromEnv()`; make `teiSubBatchSize` and the request timeouts configurable (env or handler options).

### AN-F — `TieredCache.GetOrLoad` and TEI retry/circuit not exposed as options (Tools as Primitives)
**Where:** `api/cache/tiered.go` (`GetOrLoad` bundles cache-aside + singleflight + backfill-both-tiers; `GetIfPresent`/`Set`/`GetManyIfPresent` are already proper primitives added later); `api/embed/tei.go` (retry + circuit-breaker policy in struct fields with no public options/setters).
**Why it matters:** a caller wanting different backfill policy, no singleflight, or "embed + fail fast, let me retry" has no escape hatch. The partial TieredCache refactor shows the team already moving this way.
**Fix:** finish it — promote `GetIfPresent`/`Set` as the primary surface, make `GetOrLoad`/singleflight opt-in; add `TEIClientOptions{Retry…, Circuit…, DisableCircuitBreaker}` to a `NewTEIClientWithOptions`.

---

## LOW / cleanup

### AN-G — Dead `publicMetrics` constant (and a correction to audit #3)
**Where:** `api/middleware/auth.go:19` — `publicMetrics = "/metrics"` is declared but **never referenced** in `isAuthPublicPath` (lines 62–67 list only live/ready/health/v1-healthcheck/openapi/docs).
**Note:** this means `/metrics` is actually **bearer-protected** today — so audit finding **#7 ("/metrics leaks telemetry unconditionally") is no longer true**; the constant is just dead code. (Prior re-checks conflated "constant declared" with "path public" — corrected.)
**Fix:** delete the unused `publicMetrics` constant (trivial), or — if `/metrics` is intended to be scrape-public — wire it back into `isAuthPublicPath` deliberately.

### AN-H — CORS does not expose embedding headers (Capability Discovery)
**Where:** `api/middleware/cors.go` — sets `Allow-Headers` but no `Access-Control-Expose-Headers` / `Max-Age`.
**Fix:** add `Access-Control-Expose-Headers: X-Model-Id,X-Dimensions,X-Request-Id,ETag` and `Access-Control-Max-Age: 600` so browser clients can read embedding metadata. Low (browser-facing, not agent-core).

### AN-I — Document the poll contract + single-tenant constraint (UI-Integration NA substitute)
**Where:** `README.md`.
**Fix:** state explicitly which endpoints a client polls for state transitions (warmup→`/ready`/`/warmup`; no push/SSE), and the single-bearer / single-tenant deployment constraint. Sets correct expectations; no code change.

---

## Score context (for prioritization)

| Principle | Score | Driver of which item |
|-----------|-------|----------------------|
| Action Parity | 16/31 (52%) | AN-A, AN-D |
| Tools as Primitives | 13/23 (57%) | AN-E, AN-F |
| Context Injection | 13.5/18 (75%) | AN-B, AN-C |
| Shared Workspace | 7/7 (100%) | AN-D (pre-multi-client) |
| CRUD Completeness | 2/5 (40%) | AN-A |
| Capability Discovery | 6.5/7 (93%) | AN-H |
| UI Integration | NA | AN-I (doc substitute) |
| Prompt-Native | NA | — |

The single most impactful item is **AN-A (cache flush)** — it lifts Action Parity, CRUD, and closes the orphaned-`Delete` finding in one change. This issue documents the gaps only; implementation to be PR'd separately.
