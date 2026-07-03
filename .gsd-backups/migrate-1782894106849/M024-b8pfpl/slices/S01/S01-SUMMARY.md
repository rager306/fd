---
id: S01
parent: M024-b8pfpl
milestone: M024-b8pfpl
provides:
  - Packaged performance evidence for S02 outcome decision.
requires:
  []
affects:
  - S02 performance outcome and guardrail closure
key_files:
  - benchmark-results/fd-benchmark-m024-onnx-docker1024.txt
key_decisions:
  - Use `docker restart fd-onnx-m024-bench` for the benchmark L2 restart check.
  - Use `m024-onnx-docker-benchmark` cache namespace.
  - Keep TEI default unchanged; benchmark targets only packaged ONNX image on port 18000.
patterns_established:
  - Use container-specific restart command for packaged benchmarks.
  - Validate benchmark artifact markers and raw input exclusion before closing.
  - Record metadata limitations rather than hiding them.
observability_surfaces:
  - Benchmark artifact with sanitized effective config, Redis deltas, runtime labels, artifact metadata, and summary metrics.
drill_down_paths:
  - .gsd/milestones/M024-b8pfpl/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M024-b8pfpl/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M024-b8pfpl/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T11:27:58.058Z
blocker_discovered: false
---

# S01: Packaged ONNX performance benchmark

**S01 proved the packaged ONNX Docker image remains locally performance-viable.**

## What Happened

S01 ran the packaged ONNX performance benchmark. The packaged ONNX container started, health and smoke checks passed, and the restart command targeted the packaged container. The benchmark completed with best cold latency 7.6ms, warm latency mean 2.03ms, max throughput about 864 req/s, Redis L2 restart 3.36ms, batch L1 p95 8.91ms, batch L2 p95 5.47ms, and chunk reuse warm p95 5.24ms. The artifact includes sanitized config and no raw synthetic inputs. Runtime cleanup and binary hygiene passed.

## Verification

Benchmark, artifact hygiene, cleanup, and binary hygiene passed.

## Requirements Advanced

- onnx-packaged-performance — Measured packaged ONNX Docker performance with sanitized config snapshot.

## Requirements Validated

- m024-packaged-benchmark — `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt` completed with max throughput ~864 req/s and warm latency mean 2.03ms.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None. The benchmark passed. Metadata limitation: host-side benchmark cannot hash the in-container ONNX Runtime library path.

## Known Limitations

Host-side metadata records `ONNX_RUNTIME_LIBRARY=/opt/onnxruntime/...` but cannot hash that in-container file. Future metadata can add Docker image ID/digest or in-container hash.

## Follow-ups

S02 should summarize packaged ONNX performance viability against prior TEI baseline and local ONNX evidence, while preserving experimental status.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt` — Packaged ONNX Docker benchmark artifact.
