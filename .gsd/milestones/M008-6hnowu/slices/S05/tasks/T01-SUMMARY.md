---
id: T01
parent: S05
milestone: M008-6hnowu
key_files:
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/cache/redis.go
  - api/cache/tiered.go
  - benchmark.py
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T16:41:33.629Z
blocker_discovered: false
---

# T01: Mapped fd language-sensitive layers: rewrite ROI is only plausible for proven CPU/serialization/native-inference hot paths, not generic API orchestration.

**Mapped fd language-sensitive layers: rewrite ROI is only plausible for proven CPU/serialization/native-inference hot paths, not generic API orchestration.**

## What Happened

Mapped where Go versus C/Rust could matter in fd. The current API service is mostly orchestration: Gin HTTP handlers parse JSON, iterate inputs, call a cache abstraction, encode base64/JSON responses, and use TEI as an external inference service. Redis L2 currently performs single-key GET/SET operations with binary float32 payloads; batch handlers loop per input, making round trips and serialization more likely bottlenecks than Go runtime overhead. Language choice could matter most in CPU-heavy local operations such as JSON encoding of float arrays, base64 encoding, float32 marshal/unmarshal, chunk preprocessing, or native ONNX inference if it is moved into-process. Language choice matters least when requests are cold and dominated by TEI inference latency or when cache-hit throughput is bounded by Redis/network round trips. Therefore any Go-to-Rust/C rewrite must be gated by pprof/per-layer evidence: CPU, allocation, handler timing, Redis operation timing, TEI timing, and response serialization timing.

## Verification

Read current handler/cache/benchmark code and compared against prior runtime findings.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Read: api/handlers/embeddings.go` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read: api/handlers/batch.go` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read: api/cache/redis.go` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read: api/cache/tiered.go` | -1 | unknown (coerced from string) | 0ms |
| 5 | `Read: benchmark.py` | -1 | unknown (coerced from string) | 0ms |
| 6 | `Prior benchmark context: warm cached requests around low milliseconds; cold requests dominated by model inference.` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

Current benchmark does not yet provide CPU profile, allocation profile, Redis pool stats, pprof traces, or per-layer timings. Without those, expected language-rewrite gains must remain hypotheses.

## Files Created/Modified

- `api/handlers/embeddings.go`
- `api/handlers/batch.go`
- `api/cache/redis.go`
- `api/cache/tiered.go`
- `benchmark.py`
