---
id: S01
parent: M017-j10hmp
milestone: M017-j10hmp
provides:
  - Measured S01 artifact for S02 quality outcome decision.
requires:
  []
affects:
  - S02
key_files:
  - benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt
key_decisions:
  - 512-token Go ONNX is not enough to pass strict cosine gate because very long fragments still exceed 512 tokens.
  - Ranking parity is strong enough to justify continuing ONNX remediation, but only behind experimental opt-in and with long-text policy before promotion.
patterns_established:
  - Treat evaluator exit code 2 as a measured quality FAIL when the artifact is written successfully.
  - Always use isolated Redis namespace for TEI-vs-ONNX comparisons.
observability_surfaces:
  - Sanitized legal retrieval artifact with config, selected corpus stats, ranking metrics, cross-backend cosine summaries, and worst IDs/hashes.
drill_down_paths:
  - .gsd/milestones/M017-j10hmp/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M017-j10hmp/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M017-j10hmp/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M017-j10hmp/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T07:28:39.970Z
blocker_discovered: false
---

# S01: Tagged ONNX 512 legal quality gate

**S01 measured tagged Go ONNX 512 on the full legal gate: ranking parity is strong, strict cosine still fails on >512-token fragments.**

## What Happened

S01 ran the full Russian/legal retrieval evaluator against TEI and tagged Go ONNX configured with `ONNX_MAX_SEQUENCE_LENGTH=512` and isolated namespace `m017-onnx-512-legal-quality`. The result is a measured strict FAIL: minimum cross-backend cosine is 0.98982302, below 0.999. However, ranking parity improved dramatically: top-1 agreement is 1.0, mean overlap@5 is 0.997701, and ONNX recall ratio is 1.0. The worst remaining cases are long fragments above 512 tokens, confirming that 512 is necessary but not sufficient.

## Verification

S01 verification passed as a measurement slice: service health ok, evaluator artifact written, artifact hygiene passed, and runtime cleanup verified.

## Requirements Advanced

- onnx-long-text-quality — Validated the actual tagged Go ONNX 512 runtime against the legal corpus and narrowed remaining failure to long fragments beyond 512 tokens.

## Requirements Validated

- m017-s01-onnx512-gate — `benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt` records strict FAIL with top1_agreement 1.0, mean_overlap_at_5 0.997701, recall ratio 1.0, and minimum cosine 0.98982302.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

None.

## Known Limitations

The gate still fails the strict `min_cross_backend_cosine=0.999` threshold with minimum overall cosine 0.98982302. This slice does not implement chunking.

## Follow-ups

S02 should decide that 512-token ONNX is a major improvement but insufficient for strict legal equivalence. Next remediation should be deterministic chunking or longer-sequence handling for >512-token legal fragments, followed by a full legal gate rerun.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt` — Full legal retrieval quality artifact for tagged Go ONNX at max sequence length 512.
