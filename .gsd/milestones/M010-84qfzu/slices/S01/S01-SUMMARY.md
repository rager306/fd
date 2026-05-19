---
id: S01
parent: M010-84qfzu
milestone: M010-84qfzu
provides:
  - Ranked ONNX candidate/export paths for S02/S03.
  - Required artifact provenance and output metadata fields.
  - Explicit no-production-runtime-change boundary.
requires:
  []
affects:
  - S02
  - S03
  - S04
key_files:
  - .gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md
key_decisions:
  - Primary path is model-preserving FP32 dense-only export from exact `deepvk/USER-bge-m3` snapshot.
  - Community BAAI ONNX artifacts are reference-only.
  - INT8 USER-BGE-M3 artifact is deferred until FP32 comparator/export baseline exists.
  - Comparator-before-export is the required build order.
patterns_established:
  - Comparator before runtime experiment.
  - Community ONNX artifacts are reference material unless model/revision/provenance match exactly.
  - Large model/ONNX artifacts must remain untracked and hash-recorded.
observability_surfaces:
  - ONNX provenance checklist for future benchmark/config snapshots: source, revision, export command, package versions, opset, dtype, provider, file hashes, input/output names and shapes, pooling, normalization, max sequence length, tokenizer hash, and comparator metrics.
drill_down_paths:
  - .gsd/milestones/M010-84qfzu/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M010-84qfzu/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M010-84qfzu/slices/S01/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T18:33:34.983Z
blocker_discovered: false
---

# S01: ONNX artifact provenance

**S01 established the safe ONNX provenance path: export FP32 dense-only from the exact deepvk snapshot after building a TEI comparator; do not adopt community BAAI or INT8 artifacts.**

## What Happened

S01 inspected local TEI artifact storage, researched current ONNX candidates and export paths, then synthesized a slice-level research artifact. Local storage and upstream model tree both show safetensors/tokenizer/config files for `deepvk/USER-bge-m3` but no ONNX artifact. External candidates were ranked: exact local FP32 export is the only suitable primary path; `aapot/bge-m3-onnx`, `yuniko-software/bge-m3-onnx`, and `hotchpotch/vespa-onnx-BAAI-bge-m3-only-dense` are implementation references because they use BAAI or non-FP32 variants; `skatzR/USER-BGE-M3-ONNX-INT8` is model-family relevant but out of scope for FP32. The final research artifact defines artifact hashes, output metadata, pooling/normalization checks, and build order for S02/S03.

## Verification

S01 verification passed: all three tasks are complete; `S01-RESEARCH.md` exists; the artifact includes candidate ranking, required provenance fields, dense-output risks, and explicit no-production-runtime-change guidance.

## Requirements Advanced

None.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None. S01 stayed within research/provenance scope and made no production runtime changes.

## Known Limitations

S01 did not export or load ONNX. It did not prove semantic equivalence; it only established provenance, candidate ranking, and required metadata. No ready FP32 ONNX artifact for the exact deepvk model was found.

## Follow-ups

S02 should create a dense-output comparator using fixed Russian/legal probes against the current TEI/API baseline. S03 should attempt FP32 dense-only ONNX export/load only after comparator checks exist, using ignored local artifact storage and recording hashes/output metadata.

## Files Created/Modified

- `.gsd/milestones/M010-84qfzu/slices/S01/S01-RESEARCH.md` — S01 provenance research synthesis with candidate ranking and ONNX metadata checklist.
- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T01-SUMMARY.md` — Task summaries for local artifact inspection, ONNX candidate research, and provenance synthesis.
- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T02-SUMMARY.md` — Task summaries for local artifact inspection, ONNX candidate research, and provenance synthesis.
- `.gsd/milestones/M010-84qfzu/slices/S01/tasks/T03-SUMMARY.md` — Task summaries for local artifact inspection, ONNX candidate research, and provenance synthesis.
