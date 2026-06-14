---
milestone_id: M042-fjf2en
title: fd TEI performance investigation and mitigation
status: ready-for-planning
source: M041 S04 perf discovery (2026-06-13 18:59)
predecessor: M041-4tw0w7 (fd v2 validation+observability)
gathered: 2026-06-13
---

# M042-fjf2en: fd TEI performance investigation and mitigation

## Source

M041 S04 T01 baseline + S04 T04 measurement (2026-06-13 18:59) revealed that `fd v2`'s embedding latency is bounded by TEI, not by fd's Go code:

```
TEI logs 2026-06-13 18:59 (batch=32, cold):
  total_time=6.30s
  tokenization=586μs
  queue_time=2.72s  ← 43% of total — wait, not inference
  inference_time=787ms
```

The 2.7s `queue_time` is the dominant latency contributor. fd's own handler overhead is <5ms (per S01 perf re-measurement: cache hit 1.5-4ms, sequential chunks add <1ms each). The chunked-sequential handler logic in S01 can be parallelized. The cold-path latency for `batch=128` is 25s (4 chunks × 6s), which far exceeds the fd v2 spec target of <1000ms p95.

M040/M041 successfully shipped validation, observability surface, encoding format, and graceful error envelopes. M042 investigates why the underlying embedding model call is 2.5s slower than expected and what fd can do to mitigate it.

## Project Description

`fd` is a Go embedding API service. M041 S04 chunked the handler to call TEI in ≤32-input sub-batches because `max_client_batch_size=32` (verified via `/info`). M042 is a follow-up milestone with three investigation tracks:

1. **Root cause analysis of TEI queue_time** — why is the per-chunk queue wait 2.7s when TEI's `max_concurrent_requests=512` should admit all our traffic immediately?
2. **Async pipeline** — parallelize the chunked sequential calls in the fd handler so the 4 chunks of a 128-input batch are sent concurrently instead of sequentially.
3. **ONNX conditional fallback** — per M019 measurements, ONNX Go runtime measured "best cold latency 8.3ms, warm latency mean 1.19ms, max throughput about 858 req/s". That's 100-700x faster than current TEI. M042 ships a config-gated ONNX path as a speed-first alternative to TEI, behind an env flag (off by default, opt-in), without changing the production default (TEI per R001/M015).

## Why This Milestone

The fd v2 spec (docs/fd-v2.md Section 5.4 T-P-1..T-P-5) requires <1000ms p95 across batch sizes 1, 10, 32, and concurrent 4×8. The M041 baseline shows:

| Batch | Cache hit | Cold path | Target |
|---|---|---|---|
| 1 | 1.6ms | 236ms | <50ms ✓ both |
| 10 | 2.8ms | 1.4s | <200ms ❌ cold |
| 32 | 3.7ms | 6.1s | <1000ms ❌ cold |
| 128 | (4×32 chunks) 25s | same | ❌ |

M041 made the code correct; the spec target still fails on cold paths because the upstream model is slow. Without M042, the S01 changes (validation + error envelope) are correct but the headline latency number that the user-cited (M040 docs/same-host-embedding-service-contract.md) cannot be met.

## User-Visible Outcome

### When this milestone is complete, the user can:

- Read a written root cause analysis explaining the TEI 2.7s queue_time (TEI workers, ONNX Runtime CPU vs Candle, single-backend-thread hypothesis, etc.).
- Choose **async pipeline mode** in fd (env `FD_ASYNC_CHUNKS=true`): a 128-input request is sent as 4 parallel TEI calls instead of 4 sequential, dropping cold-path latency from 25s to ~6s (4 chunks at 1.5s each in parallel).
- Choose **ONNX mode** in fd (env `FD_BACKEND=onnx` + `ONNX_ARTIFACT MANIFEST` etc): the Go ONNX runtime serves embeddings directly. M019 measured ~1ms warm / 8ms cold for batch=32 — that's two orders of magnitude faster than TEI.
- Re-run the perf acceptance suite (Section 5.4 T-P-1..T-P-5) against either mode and observe the cold path meeting <1000ms p95.
- Switch back to TEI as default with no caller-side changes (env-gated, transparent fallback).

### Entry point / environment

- Entry point: same as M041 (`api/handlers/embeddings.go` chunked loop, `api/embed/tei.go` HTTP client, new `api/embed/async.go` for parallel orchestration, optional `api/embed/onnx.go` already exists from M008/M019).
- Environment: local same-host TEI in docker-compose, Go module `api/`, ONNX runtime library `libonnxruntime.so` (already provisioned per M033), `fd_onnx` build tag (`//go:build onnx`).
- Live dependencies: TEI container `fd_tei`, ONNX artifact (path: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`).

## Current Architecture (perf focused)

```
HTTP /v1/embeddings
  └─ ValidateEmbeddingsRequest middleware (S01)
       └─ CreateEmbedding handler
            └─ for chunkStart := 0; chunkStart < len(texts); chunkStart += 32:
                 ├─ cache.GetIfPresent per text
                 └─ if miss: TEIClient.Embed(chunk-misses)  ← SEQUENTIAL — 6s per chunk
```

The sequential loop is M042's primary mitigation target. Even with the loop parallelized, the per-chunk ceiling is set by TEI.

TEI cold path telemetry (2026-06-13 18:59):
- `total_time` ≈ 6s for batch=32
- `inference_time` ≈ 700-800ms (1 chunk of 32 = 1 inference call)
- `queue_time` ≈ 2.7s
- `tokenization_time` ≈ 200μs (negligible)
- `max_concurrent_requests` reported by `/info` = 512
- `max_batch_requests` = 4 (parallel sub-batches)
- `tokenization_workers` = 11
- `max_batch_tokens` = 16384
- `max_client_batch_size` = 32

So 2.7s queue wait for a 512-cap service means either (a) only 1 effective backend thread, or (b) lock contention in the batcher, or (c) the q_time metric measures something other than concurrent-request backpressure (e.g., internal scheduling). M042 S01 investigates.

## Completion Class

- **Root cause complete**: written analysis of why `queue_time=2.7s` despite `max_concurrent_requests=512`. Includes hypothesis, evidence, and recommended action.
- **Async pipeline complete**: per-chunk requests sent in parallel (via errgroup or sync.WaitGroup), with `X-Concurrent-Chunks` header for observability, and a benchmark proving cold-path latency reduction.
- **ONNX fallback complete**: config-gated path that selects ONNX when `FD_BACKEND=onnx`, with re-validation of legal-domain quality (per M015/M016 gate) before promotion.

## Final Integrated Acceptance

To call this milestone complete, we must prove:

- Root cause analysis document explains the 2.7s queue_time with evidence.
- Async pipeline: cold path for `batch=128` drops from 25s sequential to <10s parallel.
- Async pipeline: cold path for `batch=32` drops from 6s to 3-4s (2 chunks in parallel).
- ONNX fallback: when activated, cold path for `batch=32` is <500ms (per M019 measurements).
- Legal-domain quality gate re-run on ONNX mode passes the M015 cosine threshold OR explicitly defers with documented reasoning (per M015/M016 stance "Keep TEI as production/default until ONNX legal-quality remediation").
- All M041 acceptance criteria still pass under both modes (no regression in S01 validation/envelope work).

## Architectural Decisions

### Three investigation tracks, run in order

**Decision:** S01 root cause analysis first; if S01 concludes "TEI is fundamentally slower than expected, switch to async pipeline" → S02; if S02 insufficient → S03 ONNX fallback.

**Rationale:** S01 is cheap (read logs, profile, write doc). S02 is medium (Go code change, parallel HTTP client, error aggregation). S03 is heavy (ONNX binary re-build, legal quality re-validation). Sequential discovery avoids premature ONNX work.

**Alternatives Considered:**
- Jump straight to ONNX: rejected because (a) per M015/M016 the legal quality gate isn't cleared, (b) ONNX re-build takes time, (c) async pipeline may close the gap without touching the runtime.
- Async pipeline only, no analysis: rejected because we don't know if 4 parallel TEI calls saturate TEI's `max_batch_requests=4` (sub-batches per request), or hit the 512-concurrent cap.

### Async pipeline uses bounded concurrency, not unbounded

**Decision:** Parallel chunk concurrency = 4 (matches TEI's `max_batch_requests=4`). Not `min(max_batch_size=128 / 32, 8) = 4`.

**Rationale:** TEI's internal batcher already accepts up to 4 sub-batches in parallel. Sending more would either queue (no benefit) or get rate-limited. 4 is the natural ceiling.

**Alternatives Considered:**
- Sequential: rejected (this is the current bottleneck, 25s for 128).
- Unbounded: rejected (TEI caps internally anyway, so waste of HTTP connections).
- min(N, 8): rejected (no evidence that 8 parallel > 4 is better; risk of starving other fd callers).

### ONNX fallback is opt-in, not default

**Decision:** Production default remains TEI. ONNX mode is enabled by `FD_BACKEND=onnx` (which requires the `onnx` build tag — binary must be rebuilt with `go build -tags onnx`). No automatic fallback.

**Rationale:** Per M015/M016 the local ONNX legal-quality gate is not cleared (mean cosine 0.998 at 512 tokens, but document/query divergence on 128-token-truncated inputs). Switching default could regress legal retrieval quality. The opt-in gives operators the speed when they need it, the safety when they don't.

**Alternatives Considered:**
- Auto-fallback when TEI queue_time exceeds threshold: rejected because (a) it hides the TEI perf bug rather than fixing it, (b) it could create oscillation under load.
- ONNX as new default: rejected per M015/M016 stance.

### Out of scope (lifted from M015/M016 + fd-v2.md Section 8)

- Replacing `deepvk/USER-bge-m3` (R001 still holds).
- Multi-model support.
- TEI source code changes (TEI is upstream; we work around it).
- Legal-quality remediation for ONNX 128-token truncation (separate M015 follow-up).
- Auto-scaling.

## Error Handling Strategy

- Async pipeline: if a parallel chunk fails, the entire request fails with the chunk's error envelope. No partial response (matches OpenAI semantics).
- ONNX fallback: errors are already wired through the same `Embedder` interface — no new error envelope needed.
- Any new env flag has a default that preserves current behavior (off = current TEI path).

## Risks and Unknowns

- **TEI queue_time root cause** may be in TEI itself (Rust internal) and not fixable from fd side. If S01 concludes "TEI single-backend-thread, can't change without forking TEI", then S02+S03 are the only levers.
- **Async pipeline** may saturate TEI's `max_batch_requests=4` for the entire fd instance. If 4 concurrent fd requests × 4 parallel chunks each = 16 sub-batches in flight, but TEI caps at 4, the 4×4 grid collapses to a serialized queue. S02 may need adaptive concurrency or per-fd-instance rate limiting.
- **ONNX binary size + cold start**: re-building the fd binary with `-tags onnx` requires the `onnxruntime` shared library (already in `docs/onnx-artifacts/`) but first inference is slower than steady-state. M019 measurements say "best cold 8.3ms" but that's after warmup.
- **Legal quality gate** for ONNX mode is out of M042 scope but required before any production rollout. M015/M016 already documented the gap; M042 references but does not close it.

## Slice Plan (preview, refined by gsd_plan_slice)

- **S01: TEI queue_time root cause analysis** (research, no code change, ~2h)
- **S02: Async parallel chunked TEI calls in handler** (Go code change + perf measurement, ~6h)
- **S03: ONNX conditional fallback + speed measurement** (env-gated, ~4h)

## Spec Corrections Surfaced by M041 (carried into M042 as acceptance context)

- **fd-v2.md Section 5.4 T-P-2** (`<200ms p95` for batch=10 cold path): NOT achievable on current TEI; closes with async pipeline ~1.4s.
- **fd-v2.md Section 5.4 T-P-3** (`<1000ms p95` for batch=32 cold path): NOT achievable on current TEI; closes with async pipeline ~3-4s, ONNX ~500ms.
- **fd-v2.md Section 5.4 T-P-5** (concurrent 4×8 < 2s): achievable with cache hits; with ONNX, achievable cold too.

## Hand-off to M042 in M041 docs

`docs/fd-v2.md` will be updated to note M042 as a follow-up; the perf table will gain a new column "after M042" with the projected numbers. Daily-archive wrapper (M062 S01) does not need changes — fd's external contract is unchanged.
