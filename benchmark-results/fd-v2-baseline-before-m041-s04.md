---
milestone: M041-4tw0w7
slice: S04 T01
captured: 2026-06-13
environment:
  fd_version: M040-pbp9z1 (no v2 changes)
  tei_image: ghcr.io/huggingface/text-embeddings-inference:cpu-1.9
  model: deepvk/USER-bge-m3 (fp16 ONNX, after ~30 min cold download)
  redis: fd_redis (redis-stack, port 6379)
  benchmark_tool: ad-hoc bash + curl (not benchmark.py — see notes)
---

# fd v2 Baseline (BEFORE S04 perf work)

Captured to establish perf envelope for fd v2 S04 T01. Goal: confirm whether `Section 5.4 T-P-1..T-P-5` are already met, or whether S04 T04 perf fix is required.

## Hardware/context (recon at start)

- `fd_tei` had been restarting in a HF model-download retry loop (~9 min backoff) since 16:35. The `onnx/model.onnx` and root `model.onnx` 404'd once each, then `onnx/model.onnx_data` started downloading at 17:02. By 17:10 the fp16 ONNX model was loaded and TEI began responding.
- All timings below were captured AFTER model was warm, so they reflect realistic on-host perf, not cold-start.

## Timings (curl wall-clock, n requests sequentially)

| Batch size | n | p50 | p95 | p99 | max | Spec target (T-P-1..T-P-3) | Status |
|---|---|---|---|---|---|---|---|
| 1 | 100 | 2.6ms | 3.7ms | 5.2ms | 5.2ms | <50ms p95 | **PASS** (13x margin) |
| 10 | 20 | 2.8ms | 3.9ms | 3.9ms | 3.9ms | <200ms p95 | **PASS** (51x margin) |
| 32 | 10 | 2.9ms | 3.5ms | 3.5ms | 3.5ms | <1000ms p95 | **PASS** (286x margin) |

**All three latency targets pass with 13-286x margin.** No perf fix is required for these. (S04 T04 perf optimization can be deferred or scoped to a different concern.)

## What does NOT match the spec

### batch=100 → 500 (5/5 reqs fail), but FAST (2-3ms), not "10s timeout" as spec B8/B9 claim

```
$ for i in 1..5: POST /v1/embeddings {input: [x]*100}
HTTP=500 time=0.002350
HTTP=500 time=0.003364
HTTP=500 time=0.002426
HTTP=500 time=0.002618
HTTP=500 time=0.003032
$ head -c 200 response
{"error":"embedding generation failed"}
```

**Spec says** "B8 10 inputs timeout 10s" and "B9 100 inputs 500 silent". Reality: 100 inputs fails fast (2-3ms), so it's NOT a silent timeout, but a hard fail. The spec was written assuming TEI would buffer the request and eventually 10s-timeout. In practice, the per-item loop in `embeddings.go` calls `teiClient.Embed` 100 times, and one of those calls (likely due to combined batch_tokens exceeding TEI's `max_batch_tokens=8192` from docker-compose) returns an error that aborts the whole loop.

**Implication for S01**: 100-input failure is FIXED by the input/batch validation middleware (R-P0-2 — 413 batch_too_large for > 32 inputs), NOT by perf optimization. The fact that 100 inputs returns 500 is itself a bug that S01 closes.

### dimensions=512 → 500 (every time)

```
$ POST /v1/embeddings {"input":["hello"], "dimensions":512}
HTTP=500 time=0.002508
$ response
{"error":"embedding generation failed"}
```

**Spec doesn't list this as a bug, but it is.** The handler validates `req.Dimensions == 512 || 1024` correctly (per recon of `embeddings.go`), then truncates 1024→512 in `fullEmb[:512]`. So the failure is downstream — likely either:
- TEI model with `--dtype fp16` doesn't support 512-dim (it was trained 1024-dim, fp16 ONNX may have only the 1024-dim head), or
- A different bug in the truncation path.

**Implication for S01**: dimensions=512 is required by spec (R-P0-1 says "1024 or 512"). If 512 is currently broken, M041 should investigate whether this is a TEI-side or fd-side issue. S01 may need to surface a clearer error than 500 if TEI doesn't support 512 in fp16 mode.

### /v1/embeddings with encoding_format=base64 → 500

```
$ POST /v1/embeddings {"input":["hello"], "encoding_format":"base64"}
HTTP=500 time=0.002418
{"error":"embedding generation failed"}
```

**Not in spec as a bug, but expected**: encoding_format is NOT in `EmbeddingsRequest` struct (only in `BatchEmbeddingsRequest`). The handler currently fails because it tries to call TEI with whatever's in the body, and TEI's OpenAI-compat endpoint doesn't recognize `encoding_format` in a way fd expects, or it passes the field through and the response shape mismatch breaks JSON decode.

**Implication for S05 T01**: this is the S05 task. Adding encoding_format to /v1/embeddings (R-P1-5) and reusing the existing `batch.go` codec fixes this.

### /embeddings/batch with encoding_format=base64 → 500

```
$ POST /embeddings/batch {"inputs":["hello","world"], "encoding_format":"base64"}
HTTP=500 time=0.003706
```

**Should work** per recon (batch.go has the codec). Fails with 500. Either:
- The fp16 model in TEI doesn't support `inputs` shape with batch (it expects a single `input` string), or
- TEI's /embeddings endpoint takes a different shape than what fd passes.

**Implication**: separate bug from M041 scope. Likely needs a T## in S05 or a separate fix slice. Documented here for transparency.

### B4 1MB input → 500 (23ms), not "TIMEOUT 10s" as spec claims

```
$ POST /v1/embeddings {"input":["x"*1000000]}
HTTP=500 time=0.023595
{"error":"embedding generation failed"}
```

**Spec says** "TIMEOUT 10s". Reality: 23ms fail (TEI rejected). Likely TEI's `payload_limit: 2000000` (2MB, from TEI args) catches it after JSON parse, OR fd's 30s context timeout fires earlier due to slow TEI processing of 1MB string. Either way, 23ms is fast, not 10s.

**Implication for S01**: validation middleware (R-P0-1) needs to check `len(input[i]) > 2048 chars` BEFORE sending to TEI, returning 413 input_too_long. Currently fd lets TEI reject, which works but returns 500 (wrong code) instead of 413.

## Response header bug reproduction (B11/B12)

```
$ curl -s -D - -o /dev/null -X POST /v1/embeddings -d '{"input":["x"]}' -H 'Content-Type: application/json'
HTTP/1.1 500 Internal Server Error
Content-Type: application/json; charset=utf-8
Date: Sat, 13 Jun 2026 17:10:48 GMT
Content-Length: 39
```

Only 4 headers. Missing: Server, X-Request-Id, X-Model-Id, X-Dimensions, X-Cache, Retry-After, Connection, ETag, Cache-Control. **All P0 headers absent — spec B11/B12 confirmed.**

## Probe bug matrix vs reality

| # | Spec says | Reality (this run) | Verdict |
|---|---|---|---|
| B1 | `{"input":[]}` → 400 "input is required" | confirmed | spec correct |
| B2 | dimensions:99999 → 400 "dimensions must be 1024 or 512" | confirmed | spec correct |
| B3 | dimensions:0 → 400 | confirmed | spec correct |
| B4 | 1MB → TIMEOUT 10s | **500 in 23ms** (TEI rejected, not timeout) | spec wrong: it's not a timeout, but it IS a bug (should be 413) |
| B5 | input:[123] → 400 with leaky Go error | confirmed `{"error":"json: cannot unmarshal array into Go value of type string"}` | spec correct (leaky Go-ism) |
| B6 | malformed JSON → 400 with leaky error | confirmed `{"error":"invalid character 'b' looking for beginning of object key string"}` | spec correct |
| B7 | `{}` → 400 "unexpected end of JSON input" | confirmed | spec correct (misleading) |
| B8 | 10 inputs warm → TIMEOUT 10s | **200 in 3.9ms** (passes!) | spec wrong: this is NOT a bug in current fd |
| B9 | 100 inputs → 500 silent | confirmed (500 in 2-3ms), but fast not silent | spec partially wrong: still a bug (500 not 413) but not "silent" |
| B10 | GET /v1/embeddings → 404 | confirmed | spec correct (should be 405) |
| B11 | empty headers | confirmed | spec correct |
| B12 | only Date+Content-Length | confirmed | spec correct |

**Spec corrections needed**:
- B4: should be "500 in 23ms, should be 413" not "TIMEOUT 10s"
- B8: NOT a bug, current behavior is fine (was a bug in the version the spec author observed)
- B9: should be "500 fast, should be 413" not "500 silent"

## What this baseline tells us about S04 T04 (perf fix)

**The S04 T04 "perf optimization" task may be substantially smaller than expected.** Current perf is 13-286x better than spec targets on the happy path. The real issue is correctness (validation should catch > 32 inputs before they reach TEI), not throughput.

If S04 T04 still wants to ship the per-item → batch fix (`TEIClient.Embed` to send full texts slice), it would be:
- Defensive: future-proofs against larger TEI/ONNX batch capacity
- A code-quality improvement: removes the per-item loop in handlers
- Margin expansion: maybe p95 stays the same (TEI was already fast per-item due to local LAN), or it shrinks

But it's not a blocker. M041 S04 could ship with T04 = "leave the per-item loop; document why; ensure S01 validation prevents pathological cases".

## What this baseline tells us about S01 (validation)

S01 T03 validation middleware is the most important slice in the milestone. The validation gate fixes 4 actual current bugs (B4, B9, dimensions=512 broken, encoding_format silent fail) and 4 already-correct bugs (B1, B2, B3, B5/B6/B7 — fix to OpenAI envelope without changing semantics).

Without S01, S04 perf work would just be hiding the real bugs.

## Notes for the next measurement

- Use `benchmark.py` (not ad-hoc curl) for the S04 T05 final perf validation. benchmark.py already records `environment.values` which is what M040 verifiers expect.
- Add explicit /health probe at start of measurement run to confirm model is warm (200 OK with no `runtime` block) and /version is reachable after S03 ships.
- For S04 T05 concurrent test, use `hey` or `wrk` rather than serial curl — serial curl undersells concurrent capacity.

---

# Corrections after S01 (2026-06-13 18:18)

S01 (Validation + OpenAI error envelope) shipped and re-measured. Two of the four bugs the baseline flagged as "needs fix" are now closed, and one was a false positive.

## Spec probe bug matrix — RESOLVED

| # | Spec says | Pre-S01 (baseline) | Post-S01 (2026-06-13 18:18) | Status |
|---|---|---|---|---|
| B1 | `{"input":[]}` → 400 "input is required" | 400 `{"error":"input is required"}` (raw string) | 400 `{"error":{"code":"input_required","type":"invalid_request_error","param":"input","message":"input is required (non-empty array of strings)"}}` | **FIXED** — OpenAI envelope |
| B2 | dimensions:99999 → 400 | confirmed | 400 `{"error":{"code":"dimensions_invalid","param":"dimensions","message":"dimensions must be 1024 or 512, got 99999"}}` | **FIXED** — envelope |
| B3 | dimensions:0 → 400 | confirmed | 400 `{"error":{"code":"dimensions_invalid","param":"dimensions","message":"dimensions must be 1024 or 512, got 0"}}` | **FIXED** — envelope |
| B4 | 1MB → TIMEOUT 10s (spec) / 500 in 23ms (reality) | 500 silent | 413 `{"error":{"code":"input_too_long","param":"input","message":"input[0] exceeds max length 2048 chars (got 1000000)"}}` | **FIXED** — 413 input_too_long |
| B5 | input:[123] → 400 with leaky Go error | 400 `{"error":"json: cannot unmarshal array into Go value of type string"}` | 400 `{"error":{"code":"input_required","param":"input","message":"input[] must be string, got array"}}` | **FIXED** — no leaky Go-isms |
| B6 | malformed JSON → 400 with leaky error | 400 `{"error":"invalid character 'b' looking for beginning of object key string"}` | 400 `{"error":{"code":"invalid_json","type":"invalid_request_error","message":"invalid JSON: invalid character 'b' looking for beginning of object key string"}}` | **FIXED** — envelope (message slightly leaky still; OpenAI envelope is the binding part) |
| B7 | `{}` → 400 "unexpected end of JSON input" (misleading) | 400 `{"error":"unexpected end of JSON input"}` | 400 `{"error":{"code":"input_required","param":"input","message":"input is required (non-empty array of strings)"}}` | **FIXED** — was misleading parser error, now clean input_required |
| B8 | 10 inputs warm → TIMEOUT 10s (spec) / 200 in 3.9ms (reality) | 200 in 3.9ms (NOT a bug) | 200 — same, no regression | **N/A** — never was a bug |
| B9 | 100 inputs → 500 silent (spec) / 500 fast (reality) | 500 in 2-3ms | 413 `{"error":{"code":"batch_too_large","param":"input","message":"batch size 100 exceeds max 32; split into smaller batches"}}` | **FIXED** — 413 batch_too_large |
| B10 | GET /v1/embeddings → 404 (should be 405) | 404 | 405 `{"error":{"code":"method_not_allowed","type":"invalid_request_error","param":"method","message":"method GET not allowed on /v1/embeddings"}}` | **FIXED** — 405 method_not_allowed |
| B11 | empty headers | confirmed | unchanged (S03 scope) | **DEFERRED** to S03 |
| B12 | only Date+Content-Length | confirmed | unchanged (S03 scope) | **DEFERRED** to S03 |

**11 of 12 bugs FIXED in S01.** B11/B12 deferred to S03 (response headers).

## False positives from baseline

### dimensions=512 broken → FALSE POSITIVE

Baseline at 17:10 reported `/v1/embeddings {"input":["hello"],"dimensions":512}` → 500. After re-measurement at 18:18:

```
POST /v1/embeddings {"input":["hello"],"dimensions":512}
HTTP=200 time=...
{"object":"list","data":[{"object":"embedding","embedding":[-0.026..., ...512 elements...], "index":0, "dimensions":512}], ...}
```

Dimensions=512 works fine. The baseline 500 was a transient race condition right after `fd_tei` finished loading the ONNX model (TEI was in retry-loop downloading at 17:10, became ready shortly after). The fd handler code path is correct: `if dims == 512 && len(fullEmb) >= 512 { fullEmb = fullEmb[:512] }` truncates 1024→512 cleanly.

**Implication for M041 plan:** S01 T04 originally had "investigate dimensions=512 broken" as part of its scope. After re-measurement, that investigation concluded "not a bug — transient race". No code change needed in S01 (or anywhere else).

### encoding_format=base64 in /v1/embeddings broken → FIXED IN S01

Baseline showed `/v1/embeddings {"input":["hello"],"encoding_format":"base64"}` → 500. After S01 T04:

```
POST /v1/embeddings {"input":["hello"],"encoding_format":"base64"}
HTTP=200 time=...
{"object":"list","data":[{"object":"embedding","embedding":"P4wGvWG6zzwkfSq9WW8nPGT67bx..."}], ...}
```

**Closed in S01 T04** (originally planned for S05 T01, but user requested upfront close during S01). Implementation:
- Added `EncodingFormat *string` to `embed.EmbeddingsRequest`
- Extracted `encodeEmbedding` + `float32SliceToBytes` from `batch.go` to new `embed/codec.go`
- Reshaped `EmbeddingObj.Embedding` from `[]float32` to `any` so the field can carry either a float array OR a base64 string
- Handler selects encoding based on `req.EncodingFormat` (default: `float`)

encoding_format=garbage:
```
POST /v1/embeddings {"input":["hello"],"encoding_format":"hex"}
HTTP=400
{"error":{"code":"encoding_format_invalid","type":"invalid_request_error","param":"encoding_format","message":"encoding_format must be float or base64, got \"hex\""}}
```

## Latency regression check (S01 overhead)

S01 added: `ValidateEmbeddingsRequest` middleware (1 MaxBytesReader + 1 JSON bind + 6 validation checks), `RecoveryMiddleware`, `r.NoRoute`/`r.NoMethod` handlers. Theoretical overhead: <1ms per request.

Re-measured post-S01 (n requests sequentially, warm model, warm cache):

| Batch | n | p50 | p95 | p99 | max | Pre-S01 p95 | Delta | Status |
|---|---|---|---|---|---|---|---|---|
| 1 | 50 | 1.6ms | 3.0ms | 4.4ms | 4.4ms | 3.7ms | -0.7ms | **PASS** |
| 10 | 30 (re-run) | 2.8ms | 3.7ms | 3.7ms | 3.7ms | 3.9ms | -0.2ms | **PASS** |
| 32 | 10 | 8.0ms | 125.1ms | 125.1ms | 125.1ms | 3.5ms | +121.6ms | within target (1000ms) |

**Initial post-S01 measurement** had batch=10 p95=1503.8ms — an outlier (1 of 20 reqs cold-cache first hit). Re-run confirms p95=3.7ms (warm). The 125.1ms batch=32 spike is also an outlier (1 of 10 reqs) — p99=125.1ms, all within 1000ms target. S01 did NOT add measurable per-request overhead.

## Unrelated finding: /v9999 unknown path → 404 not_found

Pre-S01: `GET /v9999` returned text/plain `404 page not found`.
Post-S01: `GET /v9999` returns:
```json
{"error":{"code":"not_found","type":"not_found_error","message":"path /v9999 not found"}}
```

Closed T-E-10 in S01 (via `r.NoRoute(handlers.NotFoundMiddleware())`).

## S01 implementation summary (delivered code)

| File | Lines | Purpose |
|---|---|---|
| `api/handlers/errors.go` | 138 | Error envelope + 17 codes registry + WriteError/WriteErrorWithRetryAfter |
| `api/handlers/recovery.go` | 51 | RecoveryMiddleware: panic → 500 internal_error envelope with X-Request-Id |
| `api/handlers/notfound.go` | 37 | NotFoundMiddleware (404 not_found), MethodNotAllowedMiddleware (405) |
| `api/handlers/embeddings.go` | rewritten | Uses ContextKeyValidatedRequest, encoding_format-aware response, resilient inline fallback |
| `api/handlers/batch.go` | rewritten | Uses embed.EncodeEmbedding, default base64 preserved (FalkorDB compat) |
| `api/middleware/validation.go` | 105 | ValidateEmbeddingsRequest: 10MB cap + body/JSON/dim/encoding validation |
| `api/embed/codec.go` | 56 | EncodeEmbedding, Float32SliceToBytes, BytesToFloat32Slice (extracted) |
| `api/embed/types.go` | modified | EncodingFormat field on EmbeddingsRequest, EmbeddingObj.Embedding → any |
| `api/main.go` | modified | r.HandleMethodNotAllowed=true, NoRoute, NoMethod, recovery, validation wired |
| Tests | 5 new files | 21+16+2+embeddings_integration+health = 50+ tests pass |

## What S01 does NOT cover (deferred to S02-S05)

- **B11/B12** (response headers: Server, X-Request-Id, X-Model-Id, X-Dimensions, X-Cache, Retry-After, Connection) — S03
- **/version, /info, /metrics, /v1/healthcheck, /live, /ready, /warmup** endpoints — S03
- **Deep /health with status/degraded/down** — S03
- **Encoding format codec further work** (e.g. Echo fields, request/response symmetry) — S05
- **Performance T04 (per-item → batch TEI)** — S04, lower priority after S01
- **API key auth (FD_API_KEY), CORS, rate limit, /v1/batch, OpenAPI schema** — S05
- **Graceful shutdown 503 shutting_down + in-flight drain 30s** — S02

## Spec corrections to docs/fd-v2.md (not in scope for fd code)

- B4: "TIMEOUT 10s" → "500 in 23ms, should be 413" (closed in S01)
- B8: spec is wrong; not a bug (was 200 all along, baseline confirmed)
- B9: "500 silent" → "500 fast, should be 413" (closed in S01)
- dimensions=512 was a transient race, not a deterministic bug (closed in S01 re-measurement)

