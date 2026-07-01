---
id: S04
parent: M014-vjfs9f
milestone: M014-vjfs9f
provides:
  - Data-backed ONNX continuation recommendation without production switch.
requires:
  []
affects:
  []
key_files:
  - benchmark-results/fd-benchmark-m014-comparison.txt
  - .gsd/DECISIONS.md
key_decisions:
  - D010: Continue ONNX only as opt-in experimental path; do not switch production/default from TEI yet.
patterns_established:
  - Runtime benchmark recommendations must separate cold/model-bound wins from cache-dominated behavior.
  - Performance evidence alone is insufficient for production switch without packaging and quality gates.
observability_surfaces:
  - Comparison artifact with metric deltas and caveats.
  - Decision D010.
  - Final verification task summary with tests/lint/artifact hygiene evidence.
drill_down_paths:
  - .gsd/milestones/M014-vjfs9f/slices/S04/tasks/T01-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S04/tasks/T02-SUMMARY.md
  - .gsd/milestones/M014-vjfs9f/slices/S04/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T04:36:29.968Z
blocker_discovered: false
---

# S04: Benchmark synthesis and decision

**S04 synthesized M014 benchmark evidence and verified the milestone closure gates.**

## What Happened

S04 compared the TEI baseline and tagged ONNX benchmark artifacts, wrote a durable comparison report, recorded decision D010, and ran final gates. ONNX is much faster for cold/model-bound paths and slightly higher in peak throughput in this run, but cache-dominated metrics are mixed and operational packaging/startup/RSS concerns remain. The recommendation is to continue ONNX as opt-in experimental work, not switch defaults.

## Verification

All final gates passed with fresh output after the last code change.

## Requirements Advanced

- onnx-performance-evidence — Converted TEI/ONNX benchmark artifacts into a reproducible recommendation.

## Requirements Validated

- benchmark-comparison — `benchmark-results/fd-benchmark-m014-comparison.txt` includes metric deltas, caveats, and recommendation.
- final-verification — Go tests, lint, tagged tests, artifact hygiene, runtime cleanup, and GitNexus checks passed.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S04 includes snapshot v2 vs v3 caveat because S03 fixed the restart-aware benchmark harness after the TEI baseline was captured. No production/default switch was made.

## Known Limitations

Single local benchmark pair only; not repeated statistical capacity testing. ONNX correctness still fixed-probe based. Native artifacts remain local/ignored and not Docker/CI packaged.

## Follow-ups

Recommended next milestone: package tagged ONNX/native tokenizer path for reproducible Docker/CI artifact supply and run repeated tuned benchmarks; keep TEI default meanwhile. Broader Russian/legal retrieval quality evaluation remains required before production switch.

## Files Created/Modified

- `benchmark-results/fd-benchmark-m014-comparison.txt` — TEI vs tagged ONNX comparison and recommendation artifact.
- `.gsd/DECISIONS.md` — Durable recommendation decision D010.
