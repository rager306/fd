---
id: T01
parent: S02
milestone: M024-b8pfpl
key_files:
  - benchmark-results/fd-onnx-docker-performance-outcome-m024-s02.txt
key_decisions:
  - Outcome compares packaged ONNX against prior M014 TEI and M019 local ONNX evidence without claiming production readiness.
duration: 
verification_result: passed
completed_at: 2026-05-20T11:29:05.930Z
blocker_discovered: false
---

# T01: Wrote the M024 packaged performance outcome artifact.

**Wrote the M024 packaged performance outcome artifact.**

## What Happened

Wrote the packaged ONNX performance outcome artifact. It records the M024 metrics, compares them to M014 TEI and M019 local ONNX, notes caveats, and lists remaining blockers. A raw synthetic input leak check found zero leaks.

## Verification

Outcome artifact exists and excludes raw synthetic benchmark inputs.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `test ! -e benchmark-results/fd-onnx-docker-performance-outcome-m024-s02.txt before write` | 0 | ✅ pass — performance_outcome_path_new=pass | 0ms |
| 2 | `write benchmark-results/fd-onnx-docker-performance-outcome-m024-s02.txt` | 0 | ✅ pass — artifact written | 0ms |
| 3 | `gsd_exec M024 outcome raw synthetic benchmark text leak check` | 0 | ✅ pass — benchmark_raw_input_leaks=0 | 37ms |

## Deviations

None.

## Known Issues

Packaged batch L1 p95 is slower than M019 local ONNX, and in-container ONNX Runtime hash is not captured by current host-side benchmark metadata.

## Files Created/Modified

- `benchmark-results/fd-onnx-docker-performance-outcome-m024-s02.txt`
