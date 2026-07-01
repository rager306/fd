# M009-zjrq6j: M009-zjrq6j: Measured cache and benchmark foundation

**Vision:** Create the measurement and cache correctness foundation needed before ONNX/provider/language optimizations, preserving current model quality while making Redis research-cache behavior tunable and benchmark results comparable.

## Success Criteria

- Benchmark config snapshots are generated and sanitized.
- Redis cache namespace and retention are env-configurable with safe defaults and tests.
- Redis persistence and memory policy are documented and reflected in benchmark artifacts.
- Batch cache-hit benchmark evidence exists before any MGET/pipeline optimization.
- No ONNX, INT8, provider, Rust, C, or model replacement work is introduced.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, every benchmark run reports the effective safe configuration needed to compare future results.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, Redis cache retention and namespace can be tuned safely for long-lived research reuse.

- [x] **S03: S03** `risk:medium` `depends:[]`
  > After this: After this, local Redis can survive restarts for research cache reuse with documented memory/eviction policy.

- [x] **S04: S04** `risk:medium` `depends:[]`
  > After this: After this, benchmark evidence shows whether batch cache hits are Redis round-trip limited.

- [x] **S05: MGET pipeline A B** `risk:high` `depends:[S04]`
  > After this: After this, batch Redis cache hits use MGET or pipelining only if baseline evidence justifies it.

## Boundary Map

## Boundary Map

| Area | In scope | Out of scope |
|---|---|---|
| Benchmark artifacts | Sanitized effective config snapshots and Redis diagnostics | ONNX provider benchmarking |
| Cache configuration | Model-aware namespace, TTL/no-expire, safe env validation | Model replacement or sparse/ColBERT cache schema |
| Redis deployment | RDB-first persistence docs/config, maxmemory/policy recording | Redis Cluster, Redis Stack, Dragonfly/Valkey |
| Batch cache performance | Baseline and conditional MGET/pipeline A/B | Full service rewrite |
| Runtime language | Keep Go service | Rust sidecar, C service |
