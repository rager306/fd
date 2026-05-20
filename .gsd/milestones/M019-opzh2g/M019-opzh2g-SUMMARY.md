---
id: M019-opzh2g
title: "ONNX 1024 performance benchmark"
status: complete
completed_at: 2026-05-20T08:22:37.855Z
key_decisions:
  - D017: ONNX 1024 is locally performance-viable after quality pass, and the next gate is artifact contract, Docker/CI packaging, and operational validation while TEI remains default.
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m019-onnx1024.txt
  - benchmark-results/fd-onnx-1024-performance-outcome-m019-s02.txt
  - .gsd/DECISIONS.md
lessons_learned:
  - Benchmark comparability requires ONNX runtime env keys, including sequence length, in sanitized snapshots.
  - ONNX 1024 keeps strong local performance despite larger sequence length.
  - A restart helper for `go run` must kill the actual listener process, not only the parent `go run` PID.
---

# M019-opzh2g: ONNX 1024 performance benchmark

**M019 proved ONNX 1024 is locally performance-viable after legal quality pass, while keeping production promotion blocked on packaging and operations.**

## What Happened

M019 benchmarked the legal-quality-passing ONNX 1024 path. S01 prepared a restart helper targeting the actual ONNX service, extended benchmark metadata to include ONNX runtime env keys, and ran `benchmark.py` against tagged Go ONNX 1024 with isolated namespace `m019-onnx-1024-benchmark`. The final benchmark recorded best cold latency 8.3ms, warm latency mean 1.19ms, max throughput about 858 req/s, Redis L2 restart 2.02ms, batch L1 p95 4.09ms, batch L2 p95 6.70ms, and chunk reuse warm p95 6.51ms. S02 compared this against TEI and prior ONNX context, then recorded D017: ONNX 1024 is locally performance-viable but remains experimental until artifact contract, Docker/CI packaging, and operational gates pass.

## Success Criteria Results

- ONNX 1024 benchmark: PASS.
- Metrics compared to TEI/prior ONNX: PASS.
- Raw text hygiene: PASS.
- Runtime cleanup: PASS.
- Next gate explicit: PASS.

## Definition of Done Results

- ONNX 1024 benchmark was run: met.
- Isolated cache namespace was used: met (`m019-onnx-1024-benchmark`).
- Benchmark artifact includes `ONNX_MAX_SEQUENCE_LENGTH=1024`: met after allowlist update and rerun.
- Raw benchmark text was excluded from artifacts: met by hygiene checks.
- Runtime cleanup verified: met.
- Fresh tests/lint/tagged checks passed: met.
- TEI remains production/default: met.

## Requirement Outcomes

- ONNX 1024 performance: validated locally.
- Benchmark comparability: improved by recording ONNX sequence length and runtime env keys.
- ONNX production readiness: remains blocked by packaging, artifact distribution, CI, and operations.
- TEI default: preserved.

## Deviations

A first benchmark run was discarded because the sanitized config snapshot did not include ONNX sequence length. `benchmark.py` was updated to allowlist safe ONNX runtime env keys, and the benchmark was rerun to produce the final artifact.

## Follow-ups

Create the next GSD milestone for ONNX 1024 artifact contract and Docker/CI packaging: document 1024 runtime contract, supply native tokenizer and ONNX artifacts without committing binaries, verify checksums and failure modes, and rerun quality/performance in packaged environment.
