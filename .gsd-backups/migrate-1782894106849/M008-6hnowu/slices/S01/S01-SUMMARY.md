---
id: S01
parent: M008-6hnowu
milestone: M008-6hnowu
provides:
  - Verified source evidence for Go/ONNX runtime options.
  - Russian legal corpus benchmark gate for future model/quantization changes.
requires:
  []
affects:
  []
key_files: []
key_decisions:
  - Model-preserving optimization is in scope first; model replacement requires Russian legal benchmark.
  - Use `go-bge-m3-embed` as a reference/candidate, not a black-box adoption.
  - Build an fd-owned ONNX dense adapter if pursuing Go ONNX path.
patterns_established:
  - Alternative models are references only unless they pass Russian legal retrieval benchmarks.
  - Current-model dense output equivalence is the first gate for ONNX experiments.
observability_surfaces:
  - Benchmark gate requires config snapshot and retrieval metrics.
drill_down_paths:
  - .gsd/milestones/M008-6hnowu/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T17:05:14.297Z
blocker_discovered: false
---

# S01: Verify model preserving embedding runtime options

**S01 verified model-preserving runtime options and defined the Russian legal quality gate.**

## What Happened

S01 verified proposed Go embedding libraries and ONNX paths from current sources. `go-bge-m3-embed` exists and is relevant but small/maturity/artifact-risky. MiniLM Go ONNX projects are useful implementation references, not replacement candidates. `yalue/onnxruntime_go` provides a viable Go wrapper path but still requires native ONNX Runtime/model/tokenizer packaging. S01 also defined the Russian legal corpus quality gate for any model-changing or quality-risking optimization.

## Verification

All S01 tasks complete and source findings recorded.

## Requirements Advanced

- R001 — Russian/legal quality gate made concrete.
- R004 — benchmark comparability extended to retrieval quality setup.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

S01 was expanded during execution after the user clarified Russian/legal quality constraints; MiniLM-like options were reclassified as technical references only.

## Known Limitations

No actual ONNX model was downloaded or benchmarked in S01; this is source verification and gate definition.

## Follow-ups

Use S01 findings in S02 integration design and S03 final recommendation. Do not adopt MiniLM or INT8 without Russian legal corpus benchmark.

## Files Created/Modified

None.
