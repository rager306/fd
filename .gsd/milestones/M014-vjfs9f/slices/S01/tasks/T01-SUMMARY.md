---
id: T01
parent: S01
milestone: M014-vjfs9f
key_files:
  - benchmark.py
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
key_decisions:
  - Benchmark matrix should include cold single request, warm repeated request, batch cached behavior, concurrency/throughput, startup time, memory/RSS, and Redis hit/miss deltas.
  - Benchmark metadata must include runtime label, build tags, native tokenizer artifact manifest/checksum/path, ONNX artifact manifest/checksum/path, ONNX Runtime shared library path/hash, Redis namespace, git commit/dirty state, Docker config hashes, and environment baseline hash.
  - Existing `benchmark.py` already supports arbitrary `BENCHMARK_API_URL`; S01 only needs snapshot enrichment, not a full benchmark rewrite.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:10:12.970Z
blocker_discovered: false
---

# T01: Defined the M014 benchmark matrix and metadata contract for TEI vs tagged ONNX comparison.

**Defined the M014 benchmark matrix and metadata contract for TEI vs tagged ONNX comparison.**

## What Happened

Inspected the existing benchmark harness and M013 artifacts. The current harness already has a sanitized config snapshot, API URL override, Redis diagnostics, batch/cache sections, and raw benchmark text exclusion. For M014 comparability, it needs optional tagged ONNX metadata fields rather than a rewrite. The benchmark matrix should compare TEI default and tagged ONNX under documented runtime labels and isolated Redis namespaces, with startup/memory signals captured outside or alongside the benchmark artifact.

## Verification

Read `benchmark.py`, M013 cosine artifact, native tokenizer manifest, and ONNX artifact manifest. The task summary lists scenarios and required metadata fields.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read benchmark.py` | 0 | ✅ pass — existing snapshot/API/cache benchmark structure inspected | 0ms |
| 2 | `read benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt` | 0 | ✅ pass — tagged ONNX correctness baseline identified | 0ms |
| 3 | `read docs/onnx-artifacts/hf-tokenizers-linux-amd64.json and user-bge-m3-dense-fp32.json` | 0 | ✅ pass — native/ONNX metadata fields identified | 0ms |

## Deviations

None.

## Known Issues

Current benchmark.py prints built-in Russian benchmark texts in source only, not artifacts; leak checks should continue to verify artifacts do not contain raw fixed comparator probe text.

## Files Created/Modified

- `benchmark.py`
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
