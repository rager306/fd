---
id: S01
parent: M038-pmw50e
milestone: M038-pmw50e
provides:
  - Verified Go ONNX API smoke proof for S02 legal/performance expansion.
requires:
  []
affects:
  - S02
key_files:
  - benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt
key_decisions: []
patterns_established:
  - Go package tests may require package-relative artifact paths; target-runtime HTTP proof should record only hashes/metadata, not raw probe text.
observability_surfaces:
  - Outcome artifact records health runtime fields, embedding dimensions/norm, namespace, and raw text logging status.
drill_down_paths:
  - .gsd/milestones/M038-pmw50e/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M038-pmw50e/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M038-pmw50e/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M038-pmw50e/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-21T10:37:25.162Z
blocker_discovered: false
---

# S01: Go ONNX runtime smoke proof

**Produced fresh Go ONNX target-runtime smoke evidence for the current artifact.**

## What Happened

S01 produced fresh Go target-runtime smoke evidence for the current ONNX artifact. Local prerequisites passed, the live Go embedder test passed, and the actual Go API running with `onnx hf_tokenizers` returned verified ONNX health metadata plus a 1024-dimensional normalized embedding under an isolated Redis namespace. The local server was stopped and port 18000 is clean.

## Verification

S01 verification passed: live Go embedder test, Go API smoke, leak checks, GitNexus detect, background process check, and port check all passed.

## Requirements Advanced

- implicit target-runtime validation requirement — Provides first fresh Go target-runtime acceptance evidence under D035.

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

S01 proves direct Go embedder and local Go HTTP smoke only. It does not yet prove legal retrieval or performance target-runtime gates.

## Follow-ups

S02 should run legal retrieval and performance drivers against actual Go endpoints if feasible, or record truthful blockers.

## Files Created/Modified

- `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt` — Go target-runtime smoke outcome with health/embedding metadata and no raw text.
