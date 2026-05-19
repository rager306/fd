---
id: S02
parent: M010-84qfzu
milestone: M010-84qfzu
provides:
  - TEI/API dense baseline artifact for S03.
  - Comparator script reusable for regenerating baseline after config changes.
  - Explicit vector shape/norm/hash/cosine evidence.
requires:
  []
affects:
  - S03
  - S04
key_files:
  - tools/compare_dense_embeddings.py
  - benchmark-results/fd-dense-comparator-m010-s02.txt
key_decisions:
  - Use `/v1/embeddings` as the TEI baseline because it returns float arrays directly.
  - Keep comparator as a separate tool under `tools/` rather than adding ONNX logic to `benchmark.py`.
  - Exclude raw probe texts from artifacts while keeping reproducible non-sensitive probe constants in source.
patterns_established:
  - Baseline artifact before ONNX export claims.
  - Sanitized probe labels/char counts instead of raw text output.
  - Float32 vector hashes for stable embedding identity checks.
observability_surfaces:
  - `tools/compare_dense_embeddings.py` for repeatable baseline capture.
  - `benchmark-results/fd-dense-comparator-m010-s02.txt` with model/dimension/probe/norm/hash/cosine evidence.
drill_down_paths:
  - .gsd/milestones/M010-84qfzu/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M010-84qfzu/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M010-84qfzu/slices/S02/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-19T18:38:04.404Z
blocker_discovered: false
---

# S02: Dense comparator baseline

**S02 established a sanitized TEI/API dense comparator baseline that S03 can use to evaluate FP32 ONNX output.**

## What Happened

S02 created a repeatable dense-output baseline for the current TEI-backed API. It defined a sanitized comparator contract, implemented `tools/compare_dense_embeddings.py`, ran it against the live local API, and saved `benchmark-results/fd-dense-comparator-m010-s02.txt`. The artifact records configuration, five probe labels/character counts, 1024-dimensional vector checks, finite-value status, L2 norms near 1.0, deterministic float32 vector hashes, pairwise cosine similarities, usage, and PASS verdict. A validation script confirmed required sections and that raw probe texts did not leak into the artifact.

## Verification

S02 verification passed: py_compile succeeded; comparator command exited 0 and produced a PASS artifact; artifact parser confirmed required sections and no raw probe text leakage.

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

None. The slice did not introduce ONNX or change production runtime defaults.

## Known Limitations

Comparator currently captures TEI/API baseline only; it does not yet run ONNX. S03 needs an ONNX export/load path and comparison step. The artifact is point-in-time evidence and should be regenerated if model/runtime config changes.

## Follow-ups

S03 should attempt FP32 dense-only ONNX export/load and compare any candidate output against `benchmark-results/fd-dense-comparator-m010-s02.txt`. If runtime/model config changes, regenerate the baseline first.

## Files Created/Modified

- `tools/compare_dense_embeddings.py` — Reusable TEI/API dense comparator script with fixed sanitized probes, vector checks, hashes, and cosine output.
- `benchmark-results/fd-dense-comparator-m010-s02.txt` — Sanitized baseline artifact from current TEI/API runtime for S03 comparisons.
