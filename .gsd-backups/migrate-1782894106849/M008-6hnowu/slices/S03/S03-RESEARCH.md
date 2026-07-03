# S03 Research: Final Optimization Recommendation

## Inputs

- S01: model-preserving Go/ONNX options and Russian legal quality gate.
- S02: integration seams and benchmark design.
- S04: Redis reusable vector cache and batch-hit throughput plan.
- S05: Go/C/Rust rewrite strategy.
- S06: ONNX CPU provider/threading/quantization plan.

## Final ranking

### 1. First next milestone: benchmark comparability and Redis reusable cache

Recommended first because it improves evidence quality and addresses a proven product need without changing model quality.

Scope:

- Add sanitized effective config snapshots to benchmark artifacts.
- Record git commit, Docker image IDs, compose hashes, env-derived cache/runtime settings, Redis INFO memory/stats, environment artifact reference, and benchmark corpus metadata.
- Add env-configured Redis cache retention:
  - TTL mode;
  - no-expire mode;
  - model/revision/tokenizer/schema/chunking namespace fields.
- Add Redis long-lived cache deployment docs/config:
  - RDB persistence first;
  - maxmemory;
  - `allkeys-lru` default or `allkeys-lfu` for repeated chunk workloads.
- Add benchmark sections for:
  - L1 hot hit;
  - Redis L2 hit after API restart;
  - batch cached inputs;
  - Redis restart/persistence reuse.

Why first:

- Required for comparing every later optimization.
- Low/medium implementation risk.
- Aligns with user requirement for research/chunking cache reuse.
- Does not risk Russian/legal embedding quality.

### 2. Second: Redis batch-hit MGET/pipeline A/B

Recommended after baseline instrumentation exists.

Scope:

- Implement bounded `MGET` or pipelined `GET` for batch cached embeddings.
- Preserve binary dense layout.
- Compare against current per-input `GET` loop.
- Measure p50/p95/p99, RPS, Redis ops/sec, hit/miss/evictions, pool stats, and memory.

Stop criteria:

- If Redis round trips do not dominate batch cached p95, do not continue Redis complexity.
- If serialization dominates, optimize Go marshal/unmarshal path first.

### 3. Third: ONNX FP32 dense-only spike

Recommended only after benchmark/config/quality gates exist.

Scope:

- Keep TEI/Candle as baseline.
- Add/factor an embedder adapter seam.
- Test BGE-M3 ONNX FP32 dense output with default CPU EP first.
- Record model/tokenizer hashes and ORT/provider/threading settings.
- Compare dense output semantics and Russian legal retrieval metrics.

Stop criteria:

- Reject if artifact provenance is unclear.
- Reject if dense output cannot match current semantics.
- Reject if Russian legal retrieval metrics regress.
- Reject if operational packaging complexity outweighs measured gain.

### 4. Fourth: ONNX threading/provider tuning

Only after FP32 dense ONNX baseline passes.

Order:

1. CPU EP default.
2. intra-op/inter-op/threading matrix.
3. optional NUMA/cpuset after topology capture.
4. oneDNN if build path is practical.
5. OpenVINO as AMD EPYC experiment only.
6. ZenDNN only if current supported build path exists.
7. INT8 only after FP32 and Russian legal quality gate.

### 5. Deferred: Rust sidecar

Rust sidecar is plausible but not first.

Use only if:

- Go ONNX wrapper/native packaging blocks required ORT controls; or
- pprof shows Go native inference glue/allocations are material; and
- sidecar can preserve dense API and quality gate.

Do not full rewrite.

### 6. Rejected for now: full C service or model replacement

C full service has high security/maintenance risk. Use C only as narrow FFI/reference if wrappers lack required ORT controls.

Model replacement is out of scope unless a Russian legal benchmark proves quality is preserved or improved.

## Next implementation milestone proposal

Title: `Measured cache and benchmark foundation`

Goal: create the measurement and cache-correctness foundation needed before ONNX/provider/language changes.

Suggested slices:

1. **Benchmark config snapshot**
   - Emit sanitized config block in every benchmark artifact.
   - Include git/Docker/env/Redis/environment fields.
   - Verify no secrets/raw texts are printed.

2. **Env-configured cache namespace and retention**
   - Add safe env parsing and validation.
   - Include model/revision/tokenizer/schema/dimension/chunking namespace fields.
   - Add TTL/no-expire modes.
   - Preserve current defaults unless explicitly configured.

3. **Redis persistence and deployment hardening**
   - Document/configure RDB-first persistence path.
   - Record Redis maxmemory/policy.
   - Keep Redis localhost binding safe.

4. **Redis batch-hit benchmark**
   - Add benchmark sections for L1, L2 after restart, cached batch, and repeated chunk reuse.
   - Capture Redis INFO and pool stats where available.

5. **MGET/pipeline A/B implementation**
   - Only after benchmark baseline exists.
   - Implement bounded batch L2 lookup if baseline shows round-trip pressure.

## Required verification for next milestone

- `docker compose config`
- `cd api && go test ./... -short`
- pinned GolangCI-Lint command
- `uv run --python 3.13 --with requests --with redis python benchmark.py`
- benchmark parser/check that config snapshot exists and is sanitized
- Redis key namespace/TTL/no-expire tests
- Redis L2 persistence/restart reuse test if persistence config changes
- GitNexus impact before editing symbols and detect changes before commit

## Non-goals for next milestone

- ONNX runtime adoption.
- INT8 quantization.
- oneDNN/OpenVINO/ZenDNN provider work.
- Rust sidecar.
- C service.
- model replacement.
- sparse/ColBERT public API.

## Bottom line

Do the boring measurement/cache foundation first. It is the safest path, creates reusable evidence for all later experiments, and directly addresses the user's long-lived research/chunking cache requirement. ONNX/provider/language work should follow only after the benchmark and quality gates can prove that it helps without degrading Russian legal embedding quality.
