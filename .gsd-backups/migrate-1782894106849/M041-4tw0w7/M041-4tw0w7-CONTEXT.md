---
milestone_id: M041-4tw0w7
title: fd v2 — service hardening and observability
status: ready-for-planning
source_doc: docs/fd-v2.md
recon: .gsd/milestones/M041-4tw0w7/slices/S01/S01-RECON.md
baseline: benchmark-results/fd-v2-baseline-before-m041-s04.md
gathered: 2026-06-13
---

# M041-4tw0w7: fd v2 — service hardening and observability

## Source

`docs/fd-v2.md` is the single source of truth for fd v2 requirements. It defines the current observed state (12 probe bugs across 5 categories), prioritized P0/P1/P2 requirements, machine-readable error catalog, OpenAPI sketch, 45 acceptance test cases, 10 behavior scenarios, and Go-flavored architecture hints. The document is binding for fd v2 implementation.

A pre-planning recon of the actual Go code (`/root/fd/api/`) was performed and is captured at `S01-RECON.md`. A live baseline timing measurement was also captured at `benchmark-results/fd-v2-baseline-before-m041-s04.md` (p50/p95/p99 for batch=1/10/32, all PASS spec targets with 13-286x margin; batch=100 fails 500 fast not timeout; dimensions=512 broken; encoding_format not in /v1/embeddings). Several risks from the original spec are now confirmed or de-risked, and several probe-bug descriptions in the spec are now corrected by the baseline.

## Project Description

`fd` is a Go embedding API service for Russian/legal-domain workloads. M040 produced the same-host service contract and the TEI-vs-ONNX runtime recommendation. M041 implements the fd v2 requirements: it makes the service production-ready for neighboring local HTTP clients by fixing functional bugs (silent timeouts, 500 on batch overflow, leaky Go error messages), adding full observability surface (deep health, /version, /info, /metrics, response headers), supporting proper lifecycle (warmup, readiness/liveness probes, graceful shutdown), and adding a small but useful feature set (LRU cache with X-Cache header, OpenAI v2 compat fields, optional API key auth, CORS, OpenAPI schema).

## Why This Milestone

M040 closed the "which runtime, what contract" question for the same-host embedding service. The remaining product gaps are operational: caller integration is fragile because validation leaks parser errors and oversized batches produce 500s; operators have no way to tell version/model/load/throughput from a single endpoint; cold starts and SIGTERM race requests into timeouts. Without addressing these, every neighboring service that calls `fd` ships its own ad-hoc retry, error-parsing, and warmup logic.

M041 turns the v2 spec from a doc into a measurable, tested, machine-verifiable service surface.

## User-Visible Outcome

### When this milestone is complete, the user can:

- Call `fd` with predictable errors: every 4xx/5xx returns OpenAI-style envelope `{"error": {"code", "type", "param", "message"}}` with a documented code from the Section 3 catalog.
- Inspect `/health` (deep), `/live`, `/ready`, `/version`, `/info`, `/metrics` for runtime state.
- Get `X-Request-Id` echoed or generated, `X-Model-Id`, `X-Dimensions`, `X-Cache: HIT|MISS`, `Retry-After` on transient failures, `Server: fd/<version>`.
- Trust the performance baseline: 1 input < 50ms p95, 10 < 200ms p95, 32 < 1000ms p95, 100 sequential zero errors, 4×8 concurrent < 2s.
- Hit a same-shape `/v1/batch`, optionally `encoding_format=base64` (~30% bandwidth savings), `user` field for per-user rate limits, `priority` hint.
- Read `/openapi.json` and `/docs` for self-describing contract.
- Optionally set `FD_API_KEY` for bearer auth.

### Entry point / environment

- Entry point: local HTTP API on `:8000`, contract documented in `docs/same-host-embedding-service-contract.md` (from M040 S01) and `/openapi.json` (new).
- Environment: local same-host service on Ubuntu/KVM host, Go module `api/`, Redis on localhost for L2 cache (existing), ONNX or TEI backend (existing).
- Live dependencies: model in GPU/CPU, Redis on localhost.

## Current architecture (from recon, 2026-06-13)

Single-process Go service on gin v1.12:

- `api/main.go` (292 lines) — entry point. Server starts blocking, only `gin.Recovery()` middleware today. Shutdown uses 10s context (need 30s), no in-flight tracking. Signal handling for SIGTERM/SIGINT exists.
- `api/handlers/embeddings.go` (95 lines) — POST /v1/embeddings. `ShouldBindJSON` for binding, per-item embed loop calling `teiClient.Embed(ctx, []string{text})` — one HTTP request to TEI per input. This is the root cause of B8 (10 inputs timeout).
- `api/handlers/batch.go` (110 lines) — POST /embeddings/batch. Already implements `encoding_format: "float" | "base64"`. Float32 LE bytes encoded via `float32SliceToBytes`. Same per-item loop pattern.
- `api/handlers/health.go` (60 lines) — GET /health. Returns `{status:"ok", time, runtime: RuntimeHealth}`. `RuntimeHealth` struct (in same file) already has `backend`, `model`, `dimensions`, `production_default`, `cache_namespace` plus ONNX verification fields — partially satisfies R-P0-8.
- `api/handlers/constants.go` — `errorKey = "error"`. JSON key for current `{"error": "string"}` shape.
- `api/embed/types.go` — `EmbeddingsRequest` (with custom `UnmarshalJSON` accepting string OR []string input), `BatchEmbeddingsRequest` (with `encoding_format` field), response types.
- `api/embed/tei.go` (104 lines) — TEI HTTP client. **CRITICAL**: `Embed(ctx, texts)` uses only `texts[0]`, ignores the rest. This explains the per-item bottleneck.
- `api/embed/onnx.go` (212 lines) — ONNX runtime embedder behind `//go:build onnx` build tag. Uses `*ort.DynamicAdvancedSession` which natively handles batches.
- `api/cache/local.go` (130 lines) — L1 `LocalCache` (sync.Map, 10000/30s). Eviction is random (Range-then-Delete), not LRU. No hit/miss counters.
- `api/cache/tiered.go` (120 lines) — `TieredCache.GetOrLoad(ctx, key, dim, loader)` with `singleflight.Group` dedup. Returns `[]float32` (not byte slice). No metric counters.
- `api/go.mod` — Go 1.25, gin v1.12, redis v9.7, validator/v10 (via gin), golang.org/x/sync (singleflight already used).

**Key finding**: tiered L1+L2 cache is already in place. The S04 T02 task is therefore not "implement LRU" but "adjust existing cache: add X-Cache header, hit/miss counters, 24h TTL env override, proper LRU eviction, expose status via return value".

**Key finding (B8 root cause)**: per-item TEI call. S04 T04 perf fix is concrete: change `TEIClient.Embed` to send the full `texts` slice, change handler to call `Embed(ctx, texts)` once instead of N times. ONNX backend already handles batches (verify before assuming).

## Completion Class

- Contract complete means: the Section 4 OpenAPI spec is implemented as endpoints and machine-verifiable; the Section 3 error catalog is the authoritative error contract.
- Integration complete means: 45 test cases in Section 5 (T-H-1..T-H-10 happy, T-E-1..T-E-15 error, T-HDR-1..T-HDR-10 headers, T-P-1..T-P-5 performance, T-E existence) pass against a running fd v2.
- Operational complete means: Prometheus metrics scrape, deep /health, /live, /ready, graceful shutdown all behave per Section 6 scenarios.

## Final Integrated Acceptance

To call this milestone complete, we must prove:

- All 4 P0 functional bugs (B4 1MB timeout, B8 10 inputs timeout, B9 100 inputs 500, B7 misleading parser error) are fixed and covered by tests.
- All 6 missing P0 endpoints (`/version`, `/info`, `/metrics`, `/v1/healthcheck`, `/live`, `/ready`) return 200 with the documented shape.
- All P0 response headers (Server, X-Request-Id, X-Model-Id, X-Dimensions, X-Cache, Retry-After, Connection) are present and tested.
- All 16 error codes from Section 3 catalog are emitted with correct HTTP status, code, and type.
- Performance baseline (Section 5.4 T-P-1..T-P-5) holds on a warm model. S04 T04 must fix the per-item TEI call which is the documented root cause.
- 10 behavior scenarios (Section 6.3 F-1..F-10) are reproduced by integration tests.
- Backward compatibility: a v1 caller (only POST /v1/embeddings with OpenAI shape) still works (with new headers and the new error envelope on errors).

## Architectural Decisions

### One milestone, five vertical slices

**Decision:** Implement fd v2 as M041 with five slices ordered by risk and dependency, not as one giant refactor.

**Rationale:** The spec is large (45 test cases, 16 error codes, 6 new endpoints, 7 response headers) and touches the entire request pipeline. Vertical slices let us validate the foundation (validation+errors, lifecycle) before building observability and features on top.

**Alternatives Considered:**
- Single mega-slice — rejected because it has too many surfaces to verify in one UAT pass.
- Per-endpoint slices — rejected because cross-cutting concerns (headers, error envelope) span all endpoints.

### Slice order: foundation → observability → perf → features

**Decision:** S01 validation/errors, S02 lifecycle, S03 observability surface, S04 performance+cache, S05 features+OpenAPI.

**Rationale:** S01 fixes the highest-severity bugs (B4/B7/B8/B9) and is a prerequisite for honest error reporting everywhere else. S02 makes warmup explicit so S04's perf numbers are trustworthy. S03 (metrics+health) is what we use to *prove* S04 perf. S05 is additive.

**Alternatives Considered:**
- Reverse order (features first) — rejected because adding features on top of leaky error surface compounds the problem.
- S03 first — rejected because metrics without validated inputs are still measuring garbage.

### Validation BEFORE model call

**Decision:** All input validation (length, batch, dimensions, JSON shape, type) happens in middleware before the request reaches the model or cache.

**Rationale:** This is the only way to return 413/400 instead of 500 (B9 root cause) and to keep p95 latency under control. The architecture hint in Section 7.5 shows the validate-first pattern.

**Alternatives Considered:**
- Validate inside the handler — rejected because it scatters the rules and makes the model path the de-facto validator.
- Validate in cache lookup — rejected because cache key is itself a function of input shape, so validation must precede key construction.

### Error envelope is OpenAI-style, not custom

**Decision:** Use the OpenAI v1 error envelope exactly: `{"error": {"code", "type", "param", "message"}}`. Section 3 catalog is authoritative.

**Rationale:** Daily-archive wrappers (httpx) already understand this shape. The alternative (custom shape) would force every caller to re-learn a fd-specific contract.

**Alternatives Considered:**
- Keep current `{"error": "string"}` and add a v2 endpoint — rejected because the v1 path is the only working one, so we must fix it in place and accept a one-time migration cost (covered by M062 S01 wrapper update).

### Reuse existing tiered cache, do not re-implement

**Decision:** S04 T02 adjusts the existing `LocalCache` (add LRU eviction, hit/miss counters, 24h TTL env override) rather than building a new LRU from scratch. The Redis L2 namespace isolation via `EMBEDDING_CACHE_VERSION` is preserved as-is.

**Rationale:** L1+L2 cache with singleflight is already battle-tested from M008/M040. The new requirements (X-Cache header, hit/miss metrics, longer TTL) are additive.

**Alternatives Considered:**
- Build a new LRU library — rejected because it would replace tested code without clear gain and risks regressing Redis namespace isolation.
- Add only metrics and skip LRU eviction improvement — rejected because the current random eviction is poor (per `LocalCache.enforceMaxSize`).

### Extract encoding_format codec to shared package

**Decision:** S05 T01 moves `encodeEmbedding` and `float32SliceToBytes` from `batch.go` into a new `api/embed/codec.go` package and reuses them from `embeddings.go` (where `encoding_format` is currently missing).

**Rationale:** The encoding logic is already correct in `batch.go`; it just needs to be applied to `/v1/embeddings` too (R-P1-5). Duplicating it would risk drift.

### /v1/batch is additive, /embeddings/batch stays for backward compat

**Decision:** S05 T04 adds `/v1/batch` (spec shape: `{"batches": [[..]]}` → `{"batches": [[..]]}`) as a NEW endpoint. Existing `/embeddings/batch` (FalkorDB shape: `{"inputs": [...], "encoding_format": "..."}`) stays untouched for backward compat.

**Rationale:** Removing or changing `/embeddings/batch` would break FalkorDB callers. The new `/v1/batch` is for callers who want the spec-shape array-of-arrays.

**Alternatives Considered:**
- Replace `/embeddings/batch` with `/v1/batch` — rejected as breaking change with no clear migration path.
- Alias `/v1/batch` → `/embeddings/batch` — rejected because shapes don't match (single array vs array of arrays).

### Out of scope (lifted from fd-v2.md Section 8)

- Replacing `deepvk/USER-bge-m3` (R001 still holds: model replacement requires legal-domain evidence).
- Multi-model support within one service run (R009 still holds: no silent per-request fallback).
- Distributed tracing (only in-memory `/v1/traces` is in scope).
- OAuth / JWT / multi-tenant (only `FD_API_KEY` bearer is in scope).
- Auto-scaling.

## Error Handling Strategy

- Validation middleware runs first; failures are 4xx with `invalid_request_error` envelope.
- Lifecycle gates (`/ready`, model loaded, shutdown) run second; failures are 503 with `overloaded_error` and `Retry-After`.
- Auth (if `FD_API_KEY` set) runs at request entry; failure is 401 with `authentication_error`.
- Rate limit (if P2-5 enabled) runs before model; failure is 429 with `rate_limit_error` and `Retry-After`.
- Model inference has a bounded `REQUEST_TIMEOUT` (default 30s); cancellation returns 504 `request_timeout`.
- Any uncaught panic is recovered, logged with `X-Request-Id`, returned as 500 `internal_error` with the same `X-Request-Id` in the message.
- Every error increments `fd_errors_total{code=...}` and is logged structured.

## Risks and Unknowns (updated after recon)

**De-risked by recon:**
- ~~Real Go code path unverified~~ — recon complete, pipeline mapped, all key files read.
- ~~B8/B9 root cause unclear~~ — confirmed: per-item TEI HTTP call (10 inputs = 10 sequential TEI requests, exceeds 10s client timeout). S04 T04 has a concrete two-line fix.
- ~~L1/L2 cache missing~~ — exists already, will be adjusted rather than re-implemented.
- ~~encoding_format implementation effort~~ — already in `batch.go`, just extract and reuse.
- ~~/health runtime block missing~~ — `RuntimeHealth` struct already in `health.go` with backend/model/dimensions/cache_namespace. S03 T02 just needs to add `/info` and `/v1/healthcheck` endpoints that expose it.

**Remaining risks:**
- **ONNX `Embed(ctx, texts)` batch semantics** — assumed to handle arrays natively via `*ort.DynamicAdvancedSession` (it's the whole point of dynamic batching), but must be verified with a unit test before S04 T04 ships. If ONNX is per-item too, the fix is the same shape but the ONNX hot path is what gets optimized.
- **`encoding_format` in `/v1/embeddings` requires response shape change** — current `EmbeddingObj.Embedding` is `[]float32`. Base64 needs to switch to `json.RawMessage` or use a custom MarshalJSON. This is a wire-format change for the `encoding_format=base64` path that MUST be coordinated with daily-archive wrappers via M062 S01.
- **`gin.Recovery()` does not preserve X-Request-Id on panic** — the new headers middleware (S03 T05) must run before Recovery, AND Recovery must be wrapped to read X-Request-Id from the gin context and include it in the 500 body. If we forget, panic-induced 500s lose request correlation.
- **Daily-archive wrapper change is out of fd's scope but is a hard dependency** — M062 S01 (daily-archive) must land in parallel or shortly after M041 S01, otherwise v1 callers see new error envelopes without consumer-side parsing. This risk was flagged in the original spec and recon confirms it.
- **Header names case-insensitive in HTTP** — the spec shows canonical-case; the implementation should use canonical-case in responses and accept any case in requests. Easy to get wrong.
