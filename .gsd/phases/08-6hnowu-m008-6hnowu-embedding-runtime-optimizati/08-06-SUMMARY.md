---
id: S06
parent: M008-6hnowu
milestone: M008-6hnowu
provides:
  - ONNX CPU/provider/quantization benchmark matrix for S03 final recommendation.
requires:
  []
affects:
  []
key_files:
  - .gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md
  - benchmark-results/fd-environment-inxi-m008.txt
key_decisions:
  - Start ONNX exploration with dense-only FP32 default CPU EP before provider-specific or INT8 variants.
  - Treat ZenDNN as caution/verify-current-support, not a simple switch.
  - Treat sparse/ColBERT outputs as future hybrid retrieval capability, not current `/v1/embeddings` behavior.
patterns_established:
  - Provider claims must be classified by current support/build practicality before benchmarking.
  - Quantization is a quality-risking optimization, not a default requirement.
observability_surfaces:
  - S06-RESEARCH benchmark config snapshot fields
  - benchmark-results/fd-environment-inxi-m008.txt
drill_down_paths:
  - .gsd/milestones/M008-6hnowu/slices/S06/tasks/T01-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S06/tasks/T02-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S06/tasks/T03-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S06/tasks/T04-SUMMARY.md
  - .gsd/milestones/M008-6hnowu/slices/S06/tasks/T05-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T17:03:39.343Z
blocker_discovered: false
---

# S06: Research ONNX CPU acceleration and quantization

**S06 turned ONNX CPU acceleration claims into a cautious, benchmarkable plan.**

## What Happened

S06 verified ONNX Runtime CPU provider options, threading/NUMA controls, INT8 quantization feasibility, and BGE-M3 dense/sparse/ColBERT output implications. The slice concluded that default CPU EP FP32 dense-only is the first safe ONNX benchmark target, while oneDNN/OpenVINO/ZenDNN and INT8 are later experiments gated by availability, reproducibility, and Russian legal quality checks.

## Verification

S06 research artifact saved and all task summaries complete.

## Requirements Advanced

- R001 — Russian/legal quality gate reinforced for INT8 and model output changes.
- R004 — benchmark config snapshots extended with ONNX/provider fields.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

No ONNX benchmark was run during this research slice; this slice only verified sources and produced a future benchmark design.

## Follow-ups

Use S06 benchmark order in final M008 recommendation. Do not implement ONNX provider changes until a future spike with config snapshots and quality gates.

## Files Created/Modified

- `.gsd/milestones/M008-6hnowu/slices/S06/S06-RESEARCH.md` — ONNX CPU acceleration research synthesis.
- `benchmark-results/fd-environment-inxi-m008.txt` — Environment snapshot used by ONNX/CPU benchmark planning.
