# S02 Research: Integration and Benchmark Design

## Scope

Map completed research branches into the current fd code/runtime seams and define the benchmark design for the next implementation milestone.

Inputs:

- S01: model-preserving Go/ONNX runtime options and Russian legal quality gate.
- S04: Redis reusable vector cache and batch-hit throughput research.
- S05: Go/C/Rust rewrite tradeoffs.
- S06: ONNX CPU provider/threading/quantization research.

## Integration seams

### 1. API contract seam

Current endpoints should remain dense and OpenAI-compatible:

- `/v1/embeddings`
- `/embeddings/batch`

Current dense shape is represented in `api/embed/types.go` and handlers.

Do not expose BGE-M3 sparse or ColBERT outputs through existing endpoints in the next milestone. If those outputs are pursued later, they need a separate hybrid retrieval design and cache/search schema.

### 2. Embedder seam

Current runtime path:

- `api/embed/tei.go` implements TEI HTTP client behavior.
- Handlers depend on dense embedding responses, not a formal multi-backend interface yet.

Next implementation should introduce or formalize an embedder adapter seam only when needed:

- `TEI` adapter remains baseline.
- `ONNX dense` adapter can be added behind the same dense vector contract.
- The first ONNX adapter must consume dense output only and compare against TEI/Candle baseline.

Do not start with Rust/C rewrite. If Go ONNX packaging works, use Go first. If Go wrapper blocks required ORT controls or pprof shows native inference glue overhead, compare a Rust sidecar later.

### 3. Cache seam

Current runtime path:

- `api/cache/redis.go` stores binary dense embeddings in Redis.
- `api/cache/tiered.go` handles L1 -> L2 -> loader cache-aside flow.
- Batch handler loops input-by-input, causing N Redis accesses on L2 hit workloads.

Next implementation should add in order:

1. model/version-aware namespace fields;
2. env-configurable TTL/no-expire mode;
3. Redis persistence/maxmemory/policy docs/config;
4. benchmark visibility for Redis INFO and cache hit/miss/pool behavior;
5. batch L2 `MGET` or bounded pipelined `GET` only after benchmark captures current baseline.

### 4. Docker/runtime seam

Runtime config should include:

- Redis persistence mode and volume expectations;
- Redis `maxmemory` and `maxmemory-policy`;
- cache TTL/no-expire mode;
- model/cache namespace settings;
- future ONNX provider/model/tokenizer/threading settings;
- safe defaults and validation.

Redis must remain bound safely as already fixed: host binding should not expose Redis externally.

### 5. Benchmark seam

`benchmark.py` should become the evidence collector for both performance and comparability.

Current gaps:

- no full sanitized effective config snapshot;
- no per-layer timing split;
- no isolated Redis round-trip/pipeline benchmark;
- no Russian legal retrieval quality metrics;
- no ONNX provider/threading/quantization matrix.

## Benchmark matrix for next implementation milestone

### Phase A: comparable baseline

Goal: make existing TEI+Go+Redis baseline reproducible.

Capture:

- git commit/branch/dirty flag;
- Docker image IDs and compose hashes;
- API/cache env settings;
- Redis INFO memory/stats/hit/miss/evictions/dbsize;
- environment artifact reference and refreshed `inxi`/`lscpu` if needed;
- benchmark corpus/input version;
- p50/p95/p99, RPS, error count.

### Phase B: Redis reusable cache and batch-hit baseline

Goal: quantify cache-hit path before optimizing it.

Scenarios:

- L1 hot hit;
- Redis L2 hit after API restart;
- Redis L2 hit with L1 disabled/reset if available;
- batch of repeated cached inputs;
- batch of unique cached inputs;
- batch miss/model call.

Metrics:

- p50/p95/p99 latency;
- RPS;
- Redis ops/sec;
- keyspace hits/misses;
- evictions;
- memory per vector/key count;
- go-redis pool wait/count if instrumented.

Stop/next condition:

- If Redis round trips dominate batch cached p95, implement MGET/pipeline.
- If serialization/unmarshal dominates, optimize Go binary/marshal path before Redis infrastructure changes.

### Phase C: Redis MGET/pipeline A/B

Goal: prove batch cache-hit improvement.

Compare:

- existing per-input GET loop;
- bounded `MGET`;
- bounded pipeline if MGET shape is not enough.

Acceptance:

- clear p95/RPS improvement on batch cached workload;
- no cache correctness regression;
- no raw text/secrets in logs/snapshots;
- memory and pool stats remain healthy.

### Phase D: ONNX FP32 dense-only baseline

Goal: test model-preserving native runtime without quality-risking shortcuts.

Compare:

- current TEI/Candle baseline;
- BGE-M3 ONNX FP32 dense output with default CPU EP;
- same input dimensions and dense output contract.

Capture:

- model/tokenizer file paths and SHA256;
- ORT version/provider;
- graph optimization level;
- thread settings;
- Docker/native library versions.

Quality gate:

- dense output semantic equivalence checks;
- Russian legal retrieval metrics before adoption.

### Phase E: ORT tuning/provider matrix

Only after Phase D passes.

Order:

1. CPU EP default;
2. intra-op/inter-op/threading matrix;
3. optional NUMA/cpuset after topology capture;
4. oneDNN if build path is practical;
5. OpenVINO as experiment;
6. ZenDNN only if current supported build path exists;
7. INT8 only after FP32 and Russian legal quality gate.

### Phase F: language sidecar experiment

Only if evidence warrants it.

Conditions:

- Go ONNX wrapper lacks required features or packaging becomes worse than Rust;
- pprof shows Go native inference glue/allocations are material;
- Rust sidecar can preserve dense API contract and quality gate.

C remains FFI/reference only.

## Required quality gate

For model-changing or quality-risking variants:

- Russian legal corpus with stable document/chunk IDs;
- query set in Russian legal language;
- relevance judgments/gold support passages;
- metrics: Recall@k, nDCG@k, MRR@k, optionally MAP/Precision@k;
- output similarity/equivalence against current baseline;
- config snapshot including corpus version and artifact hashes.

## Next implementation risk classes

Low risk:

- benchmark config snapshots;
- benchmark Redis INFO capture;
- docs for env/config;
- environment fingerprint capture.

Medium risk:

- env-driven cache TTL/no-expire and namespace validation;
- Redis persistence/maxmemory compose changes;
- batch MGET/pipeline implementation.

High risk:

- ONNX adapter/runtime packaging;
- provider-specific ORT builds;
- INT8 quantization;
- Rust sidecar.

Out of scope for next implementation:

- full Rust rewrite;
- full C service;
- model replacement;
- sparse/ColBERT public API;
- Redis Cluster/Redis Stack/Dragonfly unless single-node Redis is proven bottleneck.
