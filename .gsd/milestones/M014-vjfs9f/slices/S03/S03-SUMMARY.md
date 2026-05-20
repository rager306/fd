---
id: S03
parent: M014-vjfs9f
milestone: M014-vjfs9f
provides:
  - Tagged ONNX benchmark artifact for S04 comparison.
requires:
  - slice: S02
    provides: TEI baseline artifact for comparison.
affects:
  - S04
key_files:
  - benchmark.py
  - benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt
key_decisions:
  - Benchmark.py must use configurable `BENCHMARK_API_RESTART_COMMAND` for non-Compose benchmark targets.
  - Snapshot v3 is the restart-aware snapshot format.
  - Tokenizer-related env keys should not be omitted as secret-like solely because they contain `TOKENIZER`; real token/secret keys remain omitted.
patterns_established:
  - Non-Compose benchmark targets must provide their own restart command for L2 cache checks.
  - Tokenizer env names are not secrets by name alone; redaction should avoid false positives while preserving real secret protection.
observability_surfaces:
  - Snapshot v3 runtime metadata.
  - Configurable API restart command.
  - ONNX benchmark artifact with artifact and ORT checksums.
  - Task startup/RSS evidence for tagged ONNX server.
drill_down_paths:
  - .gsd/milestones/M014-vjfs9f/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S03/tasks/T03-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S03/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T04:31:42.576Z
blocker_discovered: false
---

# S03: Tagged ONNX benchmark

**S03 captured the tagged ONNX benchmark with corrected restart-aware metadata and cleaned up the local server.**

## What Happened

S03 validated tagged ONNX prerequisites, started the tagged server, fixed the benchmark harness so Redis L2 restart targets the benchmarked API, and ran the full ONNX benchmark with snapshot v3 metadata. The final artifact records build tag `hf_tokenizers`, native tokenizer checksum, ONNX artifact checksum, ONNX Runtime library checksum, isolated Redis namespace, and benchmark results. Artifact hygiene and raw-text leak checks passed. The local ONNX server and transient runtime files were removed afterward.

## Verification

Tagged preflight, server health, restart-aware benchmark run, artifact parser/leak checks, correctness-gate reference, cleanup, and GitNexus scope review passed.

## Requirements Advanced

- onnx-performance-evidence — Produced the tagged ONNX performance artifact under isolated cache namespace and tagged/native metadata.
- benchmark-comparability — Fixed restart metadata/harness behavior needed for non-default benchmark comparability.

## Requirements Validated

- tagged-onnx-benchmark — `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt` includes snapshot_version 3, tagged metadata, checksums, and benchmark sections; hygiene checks passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S03 discovered and fixed a benchmark harness issue: Redis L2 restart must target the benchmarked ONNX server, not always Docker Compose `api`. Snapshot_version advanced from 2 to 3 for restart-aware artifacts. Also fixed a safe-env redaction false positive where `TOKENIZER` was treated as a secret token.

## Known Limitations

Tagged ONNX benchmark used local ignored native/ONNX artifacts and a transient local restart helper; this is not Docker/CI production packaging. Startup/RSS was captured in task summary rather than embedded in benchmark artifact. TEI baseline artifact remains snapshot v2 while ONNX is snapshot v3 due the discovered restart fix.

## Follow-ups

S04 should compare TEI vs tagged ONNX metrics. Important ONNX summary: best cold 10.2ms, warm mean 1.63ms, max throughput ~891 req/s at 4 concurrent, Redis L2 restart 2.70ms, batch L1 p95 5.62ms, batch L2 p95 4.41ms, chunk reuse warm p95 9.00ms. Compare against TEI: best cold 59.0ms, warm mean 2.25ms, max throughput ~750 req/s at 16 concurrent, Redis L2 restart 2.82ms, batch L1 p95 4.16ms, batch L2 p95 5.51ms, chunk reuse warm p95 7.04ms.

## Files Created/Modified

- `benchmark.py` — Added configurable benchmark API restart command, snapshot v3 metadata, and tokenizer-aware safe env redaction.
- `benchmark-results/fd-benchmark-m014-onnx-hf-tokenizer.txt` — Tagged ONNX benchmark artifact with runtime/native/ONNX/ORT metadata and performance metrics.
