---
id: S02
parent: M019-opzh2g
milestone: M019-opzh2g
provides:
  - Next gate recommendation for ONNX 1024 packaging and CI validation.
requires:
  []
affects:
  - Future ONNX 1024 packaging and CI milestone
key_files:
  - benchmark-results/fd-onnx-1024-performance-outcome-m019-s02.txt
  - .gsd/DECISIONS.md
  - benchmark.py
key_decisions:
  - D017: ONNX 1024 is locally performance-viable after legal quality pass, but remains experimental pending packaging/CI/artifact/operational gates.
  - Next immediate gate is packaging/CI/artifact contract, not performance tuning.
patterns_established:
  - Benchmark snapshots must include ONNX runtime env keys for comparability.
  - Local performance viability is a gate to packaging work, not production authorization.
observability_surfaces:
  - Performance outcome artifact comparing TEI M014, ONNX M014, and ONNX 1024 M019 metrics.
drill_down_paths:
  - .gsd/milestones/M019-opzh2g/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M019-opzh2g/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M019-opzh2g/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T08:21:22.182Z
blocker_discovered: false
---

# S02: Performance outcome decision

**S02 recorded that ONNX 1024 is locally performance-viable but not production-ready.**

## What Happened

S02 assessed the ONNX 1024 benchmark and recorded D017. The measured runtime is performance-viable on the current local host: best cold latency 8.3ms, warm mean 1.19ms, max throughput about 858 req/s, Redis L2 restart 2.02ms, batch L1 p95 4.09ms, batch L2 p95 6.70ms, and chunk reuse warm p95 6.51ms. These results support moving to packaging and CI validation while keeping TEI as production/default.

## Verification

Fresh verification passed: Python compile and artifact hygiene, Go short tests, pinned lint, tagged HF tokenizer tests, clean runtime, and GitNexus low impact checks.

## Requirements Advanced

- onnx-1024-performance — Moved ONNX 1024 from performance unknown to packaging/operations blocked.

## Requirements Validated

- m019-1024-performance-decision — `benchmark-results/fd-onnx-1024-performance-outcome-m019-s02.txt` and D017 state ONNX 1024 is locally performance-viable but remains experimental.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

M019 does not validate Docker/CI reproducibility, native tokenizer distribution, ONNX artifact distribution, or production operations.

## Follow-ups

Plan a future milestone for ONNX 1024 artifact contract, Docker/CI packaging, checksum verification, native tokenizer provisioning, and operational diagnostics. Re-run quality and performance in packaged environment before any production promotion.

## Files Created/Modified

- `benchmark-results/fd-onnx-1024-performance-outcome-m019-s02.txt` — Outcome assessment for ONNX 1024 performance benchmark.
- `.gsd/DECISIONS.md` — Decision register updated with D017.
- `benchmark.py` — Benchmark sanitized env allowlist includes ONNX runtime keys for comparability.
