---
id: S01
parent: M018-vq2ttb
milestone: M018-vq2ttb
provides:
  - Measured S01 artifact for S02 outcome decision.
requires:
  []
affects:
  - S02
key_files:
  - benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt
key_decisions:
  - Tagged Go ONNX 1024 passes strict legal quality thresholds on the selected corpus gate.
  - Quality no longer blocks the 1024 path for this corpus, but production readiness still requires performance and operational gates.
patterns_established:
  - Longer sequence length can be a simpler remediation than chunking when it passes quality; performance remains a separate gate.
  - Keep quality, performance, packaging, and promotion gates separate.
observability_surfaces:
  - Sanitized legal retrieval artifact with config, selected corpus stats, ranking metrics, cross-backend cosine summaries, and worst IDs/hashes.
drill_down_paths:
  - .gsd/milestones/M018-vq2ttb/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M018-vq2ttb/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M018-vq2ttb/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M018-vq2ttb/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T07:39:17.230Z
blocker_discovered: false
---

# S01: Tagged ONNX 1024 legal quality gate

**S01 proved tagged Go ONNX 1024 passes the full legal quality gate on the selected corpus.**

## What Happened

S01 ran the full Russian/legal retrieval evaluator against TEI and tagged Go ONNX configured with `ONNX_MAX_SEQUENCE_LENGTH=1024` and isolated namespace `m018-onnx-1024-legal-quality`. The result passed strict thresholds: minimum cross-backend cosine is 0.99989883, top-1 agreement is 1.0, mean overlap@5 is 0.997701, and ONNX recall ratio is 1.0. Runtime was cleaned up after the run.

## Verification

S01 verification passed: service health ok, evaluator PASS artifact written, artifact hygiene passed, and runtime cleanup verified.

## Requirements Advanced

- onnx-long-text-quality — Validated a 1024-token tagged Go ONNX path against the strict legal corpus quality gate.

## Requirements Validated

- m018-s01-onnx1024-gate — `benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt` records PASS with minimum cosine 0.99989883, top1_agreement 1.0, mean_overlap_at_5 0.997701, and recall ratio 1.0.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

This slice validates quality only. It does not validate performance, memory pressure, Docker/CI packaging, artifact distribution, or production rollout.

## Follow-ups

S02 should record that ONNX 1024 passes quality and the next gate should be performance, memory, artifact/package, and CI validation before any production promotion. Chunking remains a future policy for documents beyond the validated corpus/sequence length, not the immediate blocker for this gate.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt` — Full legal retrieval quality artifact for tagged Go ONNX at max sequence length 1024.
