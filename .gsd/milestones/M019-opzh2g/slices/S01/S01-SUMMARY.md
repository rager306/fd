---
id: S01
parent: M019-opzh2g
milestone: M019-opzh2g
provides:
  - ONNX 1024 performance artifact for S02 decision.
requires:
  []
affects:
  - S02
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m019-onnx1024.txt
key_decisions:
  - ONNX runtime env keys must be included in benchmark snapshots for comparability.
  - ONNX 1024 benchmark is performance-viable enough to assess packaging/CI as the likely next gate, pending S02 comparison.
patterns_established:
  - Benchmark artifacts must include ONNX sequence length and runtime env, not only benchmark-prefixed variables.
  - Runtime restart helpers must kill the actual listener, not just the `go run` parent PID.
observability_surfaces:
  - Benchmark artifact with effective config snapshot, artifact checksums, Redis metadata, latency/throughput/cache metrics.
drill_down_paths:
  - .gsd/milestones/M019-opzh2g/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M019-opzh2g/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M019-opzh2g/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M019-opzh2g/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T08:17:11.310Z
blocker_discovered: false
---

# S01: ONNX 1024 performance benchmark

**S01 benchmarked ONNX 1024 and captured strong local performance metrics with sanitized config.**

## What Happened

S01 benchmarked the legal-quality-passing ONNX 1024 runtime. The final artifact records `ONNX_MAX_SEQUENCE_LENGTH=1024`, isolated namespace `m019-onnx-1024-benchmark`, native tokenizer and ONNX artifact manifests, restart command semantics, and summary metrics. ONNX 1024 achieved best cold latency 8.3ms, warm latency mean 1.19ms, max throughput about 858 req/s, Redis L2 restart 2.02ms, batch L1 p95 4.09ms, batch L2 p95 6.70ms, and chunk reuse warm p95 6.51ms. Runtime cleanup left port 18000 clean.

## Verification

S01 verification passed: benchmark exited 0, artifact contains ONNX 1024 config, artifact hygiene passed, and runtime cleanup was verified.

## Requirements Advanced

- onnx-1024-performance — Validated local performance feasibility of ONNX 1024 after legal quality pass.

## Requirements Validated

- m019-s01-onnx1024-benchmark — `benchmark-results/fd-benchmark-m019-onnx1024.txt` records ONNX 1024 config and performance summary.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

First benchmark artifact was discarded and rerun after adding ONNX runtime env keys to benchmark.py's sanitized allowlist. This improved comparability and the final artifact includes `ONNX_MAX_SEQUENCE_LENGTH=1024`.

## Known Limitations

Single local host benchmark only; no Docker/CI packaging, production rollout, or broader load testing. Benchmark helper in `.gsd/runtime` is runtime-only and not committed.

## Follow-ups

S02 should compare M019 ONNX 1024 benchmark with M014 TEI and ONNX 128/512 context, then choose whether next gate is packaging/CI or performance tuning.

## Files Created/Modified

- `benchmark.py` — Benchmark config allowlist extended to include ONNX runtime and threading env keys.
- `benchmark-results/fd-benchmark-m019-onnx1024.txt` — ONNX 1024 benchmark artifact with sanitized effective config and summary metrics.
