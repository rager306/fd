# S02: Assess integration and benchmark design

**Goal:** Map validated research branches to the current fd API/cache/benchmark/Docker integration seams and define benchmark design for the next implementation spike.
**Demo:** After this, possible integration paths are mapped against the current API and benchmark setup.

## Must-Haves

- Integration seams are mapped to current files and runtime boundaries.
- Benchmark matrix covers Redis, ONNX, language/runtime sidecar, and quality gates.
- Risks and stop criteria are explicit.

## Proof Level

- This slice proves: research synthesis from completed S01/S04/S05/S06 evidence

## Integration Closure

Produces integration/benchmark design that S03 can rank into a final recommendation.

## Verification

- Defines what timing, config, and diagnostic fields future benchmarks must capture.

## Tasks

- [x] **T01: Map fd integration seams** `est:small`
  Synthesize current fd integration seams from prior research: embedder boundary, cache boundary, handler/API contract, Docker/runtime configuration, and benchmark harness. Identify which changes are low-risk config/benchmark changes versus code-path changes needing GitNexus impact analysis.
  - Files: `api/embed/tei.go`, `api/embed/types.go`, `api/cache/redis.go`, `api/cache/tiered.go`, `api/handlers/embeddings.go`, `api/handlers/batch.go`, `benchmark.py`, `docker-compose.yaml`
  - Verify: Review against current source map and completed research artifacts.

- [x] **T02: Define benchmark matrix and stop criteria** `est:small`
  Define the benchmark matrix for the next implementation milestone: baseline, Redis long-lived cache, Redis batch-hit MGET/pipeline, ONNX FP32 dense-only CPU EP, ORT threading variants, optional provider variants, INT8 only after quality gate, and Rust sidecar only after native inference evidence. Include quality metrics, latency/resource metrics, config snapshots, and stop criteria.
  - Files: `benchmark.py`, `benchmark-results/fd-environment-inxi-m008.txt`
  - Verify: Research artifact includes metrics, config fields, and ordered experiment matrix.

## Files Likely Touched

- api/embed/tei.go
- api/embed/types.go
- api/cache/redis.go
- api/cache/tiered.go
- api/handlers/embeddings.go
- api/handlers/batch.go
- benchmark.py
- docker-compose.yaml
- benchmark-results/fd-environment-inxi-m008.txt
