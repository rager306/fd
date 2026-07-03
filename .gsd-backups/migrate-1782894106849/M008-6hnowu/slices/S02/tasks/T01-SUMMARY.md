---
id: T01
parent: S02
milestone: M008-6hnowu
key_files:
  - api/embed/tei.go
  - api/cache/redis.go
  - api/cache/tiered.go
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - benchmark.py
  - docker-compose.yaml
  - docker-compose.override.yaml
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:13:04.310Z
blocker_discovered: false
---

# T01: Mapped current integration seams: keep API dense and Go-based, add config/benchmark evidence first, then cache and ONNX adapters behind existing seams.

**Mapped current integration seams: keep API dense and Go-based, add config/benchmark evidence first, then cache and ONNX adapters behind existing seams.**

## What Happened

Mapped fd integration seams. The embedder seam is currently `TEIClient` and dense embedding types; future ONNX work should add an adapter behind the same dense contract rather than changing handlers first. The cache seam is `RedisCache` + `TieredCache`; next safe changes are config/env parsing, model-aware key namespace, TTL/no-expire retention, metrics/logging, and then batch-aware L2 lookup. Handler/API seam should stay OpenAI-compatible dense `/v1/embeddings` and `/embeddings/batch`; sparse/ColBERT remains out of scope. Docker/runtime seam should carry Redis persistence/maxmemory/policy and ONNX artifact/env settings without exposing secrets. Benchmark seam is `benchmark.py`; it should become the source of comparable evidence by recording sanitized config, environment fingerprint, Redis INFO, pool stats where available, per-layer timings, and quality gate references. Low-risk first changes: benchmark snapshots, README/config docs, env parsing with defaults. Higher-risk code-path changes: batch MGET/pipeline, ONNX adapter, threading/provider runtime knobs.

## Verification

Cross-checked completed S01/S04/S05/S06 research against current source map and known files.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Read: .gsd/milestones/M008-6hnowu/slices/S01/S01-SUMMARY.md` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read: .gsd/milestones/M008-6hnowu/slices/S04/S04-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read: .gsd/milestones/M008-6hnowu/slices/S05/S05-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read: .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

This mapping does not edit code. Before future implementation, run GitNexus impact on edited symbols such as `RedisCache`, `TieredCache`, `CreateEmbedding`, `CreateBatchEmbeddings`, `TEIClient`, or benchmark functions.

## Files Created/Modified

- `api/embed/tei.go`
- `api/cache/redis.go`
- `api/cache/tiered.go`
- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `benchmark.py`
- `docker-compose.yaml`
- `docker-compose.override.yaml`
