---
slice: S01
milestone: M041-4tw0w7
title: Current fd pipeline recon
executed: 2026-06-13 during M041 planning
---

# S01 Recon: Current fd request pipeline

## Pipeline (as observed in /root/fd/api/)

```
HTTP request
  └─► http.Server (main.go: srv := &http.Server{...})
       └─► gin.New() + gin.Recovery() (only middleware today)
            └─► Router: r.GET("/health", ...) / r.POST("/v1/embeddings", ...) / r.POST("/embeddings/batch", ...)
                 └─► Handler (gin.HandlerFunc)
                      └─► ShouldBindJSON(&req)         ← validation is here, per-handler
                      └─► for i, text := range texts:  ← per-item loop, single text per call
                           └─► tiered.GetOrLoad(ctx, text, dims, loader)
                                └─► L1 (LocalCache, 10000/30s) + L2 (Redis) + singleflight
                                     └─► loader: teiClient.Embed(ctx, []string{text})  ← PER-ITEM HTTP REQUEST
                                          └─► POST {teiURL}/embeddings
```

## Files of interest

| File | Role | Notes |
|---|---|---|
| `api/main.go` (292 lines) | entry point, server setup, signal handling | shutdown is 10s context, no in-flight tracking, no /ready, no /metrics |
| `api/handlers/embeddings.go` (95 lines) | POST /v1/embeddings | `errorKey` is the JSON key for `{"error": "..."}`; per-item loop; 30s context timeout |
| `api/handlers/batch.go` (110 lines) | POST /embeddings/batch | encoding_format ALREADY implemented (float/base64) |
| `api/handlers/health.go` (60 lines) | GET /health | returns `{status:"ok", time, runtime: RuntimeHealth}` — `runtime` block already satisfies R-P0-8 partially |
| `api/handlers/constants.go` (3 lines) | `errorKey = "error"` | reuse for envelope migration |
| `api/embed/types.go` | EmbeddingsRequest/Response, BatchEmbeddingsRequest/Response | `Dimensions *int` is pointer; `UnmarshalJSON` accepts string OR []string input |
| `api/embed/tei.go` (104 lines) | TEI HTTP client | **CRITICAL BUG**: `Embed(ctx, texts)` uses ONLY `texts[0]`, ignores the rest. This is the per-item bottleneck that explains B8 (10 inputs timeout). |
| `api/embed/onnx.go` (212 lines) | ONNX runtime embedder | build tag `//go:build onnx`. Uses `*ort.DynamicAdvancedSession`. |
| `api/cache/local.go` (130 lines) | L1 LocalCache (sync.Map, 10000/30s) | already exists, no X-Cache header on responses |
| `api/cache/tiered.go` (120 lines) | L1+L2 with singleflight | already exists, no metric counters |
| `api/go.mod` | Go 1.25, gin v1.12, redis v9.7, validator v10 | validator available via gin → can reuse for struct validation |

## What already exists (de-risks several plan items)

| Plan item | Already in code? | What's missing |
|---|---|---|
| **R-P0-4 /health deep** | Partial — `RuntimeHealth` struct exists with backend/model/dimensions/cache_namespace | status:ok/degraded/down, last_inference_at, in_flight, model_loaded, warmup_done, device |
| **R-P1-4 LRU cache** | YES — tiered L1 10000/30s + L2 Redis + singleflight | X-Cache header on responses, fd_cache_hits_total metric, 24h TTL env override, proper LRU eviction |
| **R-P1-5 encoding_format=base64** | YES — implemented in `batch.go` (base64 float32 LE encoding) | needs to be moved to `/v1/embeddings` |
| **R-P0-8 /info** | Partial — RuntimeHealth has the data | new endpoint that returns the data |
| **R-P0-7 /version** | NO | needs buildinfo package + ldflags injection |
| **R-P0-9 /metrics** | NO | needs Prometheus client |
| **R-P0-5 graceful shutdown** | Partial — 10s Shutdown(ctx), SIGTERM/SIGINT handled | needs 30s + in-flight tracking + 503 shutting_down for new requests |
| **R-P0-3 503 on model not loaded** | NO | lifecycle state doesn't exist |
| **R-P0-4 /ready, /live** | NO | needs lifecycle state + probes handler |
| **R-P0-18 OpenAI error envelope** | NO — current is `gin.H{errorKey: "string"}` | needs ErrorResponse struct + 16 codes |
| **R-P0-11..R-P0-17 response headers** | NO | headers middleware doesn't exist |
| **R-P0-1, R-P0-2 validation** | NO — no body size limit, no batch size 32 limit, no input length 512 tokens check, no per-element type check | needs validation middleware |
| **R-P0-6 perf baseline** | NO | per-item TEI call → batch fix in S04 |

## Root cause of the B8/B9 performance bugs

```
Handler.CreateEmbedding:
  for i, text := range texts {
      emb := tiered.GetOrLoad(text, dims, loader)  // each call separate
        loader: teiClient.Embed(ctx, []string{text})  // ← ONE TEXT PER HTTP REQUEST
  }
```

10 inputs = 10 sequential HTTP requests to TEI. With 10s client timeout and TEI needing ~1s per request, B8 (timeout) is expected. The fix is in S04 T04: change the loop to make a single `Embed(ctx, texts)` call, AND change `TEIClient.Embed` to send `Input: texts` (full array) instead of `texts[0]`. `ONNXEmbedder` already handles batches via `*ort.DynamicAdvancedSession` (verify before assuming).

## Critical gotchas to encode in plan

1. **`/health` is NOT currently returning model_loaded/warmup_done/in_flight** — must be added in S03, not assumed from M040.
2. **`EMBEDDING_CACHE_VERSION` namespace** — already in use via `redisOptions.Namespace.String()` from env. S04 must NOT regress this: keep Redis L2 namespacing.
3. **TEI URL injection** — `TEI_URL` env var. S04 perf work should not change `TEI_URL` semantics.
4. **`Benchmark.py` reads `environment.values`** (nested map) — any new S04 perf artifact must put env values in that map for benchmark compatibility.
5. **ONNX build tag `//go:build onnx`** — any new ONNX code path must keep this tag; the `onnx_disabled.go` file likely has the default stub.
6. **`validator/v10`** is in deps — use for `binding:"required,max=32,min=1,dive,max=2048"` style validation instead of manual loops.
7. **`gin.Recovery()` is in place** — panics currently give a default 500 without our envelope; S01 T04 must wrap Recovery to return OpenAI envelope.
8. **No `api/middleware/` dir** — middleware/ subpackages will be new files. No precedent in repo.

## Required changes to S01..S05 plans

### S01 (Validation + envelope)

- T01 recon: this file. Task is now "formally write S01-RECON.md into milestone dir" — 30min, not 2h.
- T03 validation middleware: use `binding:"required,max=32,dive,max=2048"` struct tags + `c.Request.Body = http.MaxBytesReader(w, body, 10*1024*1024)` for size limit. Manual loop only for cross-field checks.
- T04 wire: also wrap `gin.Recovery()` to return OpenAI envelope on panic with X-Request-Id. Replace per-handler `gin.H{errorKey:...}` with `errors.NewError(...).WriteJSON(c)`.

### S02 (Lifecycle)

- T05 graceful shutdown: change `context.WithTimeout(ctx, 10*time.Second)` → 30s. Add `sync.WaitGroup` for in-flight, set `lifecycle.State.BeginShutdown()` BEFORE `srv.Shutdown()`. New requests during shutdown must get 503 via lifecycle middleware (T04) before reaching handler.

### S03 (Observability surface)

- T01 buildinfo: ldflags injection via `//go:linkname` or `var` package vars. Set defaults so dev builds still work.
- T02 /info: REUSE `RuntimeHealth` as the data source; wrap in `{models: [runtime]}` array shape per spec.
- T02 /v1/healthcheck: alias handler that calls `writeHealth` with same `RuntimeHealth`. Trivial.
- T05 headers middleware: must run BEFORE validation/lifecycle so even 4xx/5xx have `X-Request-Id` and `Server` set. Recovery wrapper also needs to preserve X-Request-Id on panic path.

### S04 (Perf + cache)

- T01 baseline: USE existing `benchmark.py` for p50/p95/p99 measurement rather than rolling a new `measure_fd_baseline.sh`. benchmark.py already records `environment.values` for artifact compatibility.
- T02 LRU cache: do NOT implement new LRU. ADJUST existing `LocalCache`: add LRU eviction (currently random via `data.Range`), expose `HitCount/MissCount` counters, add `FD_CACHE_TTL_HOURS` env (default 24h), pass-through X-Cache status via callback/return value.
- T03 cache middleware: small new file that wraps `TieredCache.GetOrLoad` calls and sets `X-Cache: HIT|MISS` via `c.Header(...)`. Sit in handler (not engine middleware) because tiered cache is per-call.
- T04 perf optimization: TWO real changes:
  1. `TEIClient.Embed`: change `text := texts[0]` → marshal `Input: texts` (the full slice). TEI supports batch input natively. Verify against current TEI.
  2. `EmbeddingsHandler.CreateEmbedding`: change per-item loop to single `h.teiClient.Embed(ctx, texts)` call. Apply cache per-item (still, because cache key is per-text) but model call is one batch.
  - Verify `ONNXEmbedder.Embed` already handles batches (it should, given `*ort.DynamicAdvancedSession`).
- T05 perf validation: extend `benchmark.py` with fd v2 specific cases rather than new script. Save to `benchmark-results/fd-v2-perf-validation-m041-s04.md`.

### S05 (OpenAI v2 + P2)

- T01 encoding_format: EXTRACT `encodeEmbedding` and `float32SliceToBytes` from `batch.go` into `api/embed/codec.go` (new), import in `embeddings.go`. Add `EncodingFormat *string` to `EmbeddingsRequest`. Validation rejects non-float/non-base64.
- T04 /v1/batch: ADDITIVE new endpoint at `/v1/batch` with `{"batches": [[..]]}` shape. Keep `/embeddings/batch` for backward compat. Per-inner-batch process through validation+lifecycle+cache+model pipeline.
- T07 OpenAPI: generate spec from Go types via `kin-openapi` or hand-written. Include all existing routes (/health, /v1/embeddings, /embeddings/batch) PLUS new v2 routes.

## Migration impact (M062 daily-archive)

- `src/arxiv_archive/embedder.py` (httpx) and `scripts/m057_*.py` (urllib) parse `{"object": "list", "data": [...]}`. The new error envelope `{"error": {"code": ...}}` will break them ONLY on errors, not on success. M062 S01 must update these scripts to check `response.get("error", {}).get("code")` first.
- The new headers (X-Request-Id, X-Model-Id, X-Dimensions) are additive and don't require caller changes.
- The new `encoding_format` field is OPT-IN — callers ignoring it get the same float array as before.
- This means S01 migration risk is contained: existing happy-path callers don't break, error-path callers need M062 S01 update.
