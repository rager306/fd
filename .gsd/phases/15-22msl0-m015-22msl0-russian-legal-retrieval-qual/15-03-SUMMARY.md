---
id: S03
parent: M015-22msl0
milestone: M015-22msl0
provides:
  - Failed legal quality gate evidence for S04 closure.
requires:
  - slice: S02
    provides: Evaluator and dry-run hygiene.
affects:
  - S04
key_files:
  - benchmark-results/fd-legal-retrieval-m015-s03.txt
  - tools/evaluate_legal_retrieval.py
key_decisions:
  - Treat S03 as a failed legal quality gate for production-readiness purposes.
  - Do not proceed to ONNX packaging/tuning as the next priority without addressing long-text quality divergence.
patterns_established:
  - Quality gate artifacts must include worst-case IDs/hashes for failures.
  - A fast ONNX path is not sufficient when legal-corpus vector equivalence fails on long texts.
observability_surfaces:
  - Live quality gate artifact with worst document/query IDs, lengths, hashes, and metrics.
  - Runtime startup and cleanup evidence.
drill_down_paths:
  - .gsd/milestones/M015-22msl0/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M015-22msl0/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M015-22msl0/slices/S03/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T05:03:22.596Z
blocker_discovered: false
---

# S03: Run legal retrieval quality gate

**S03 produced a failed Russian/legal quality gate for tagged ONNX due severe cross-backend cosine outliers.**

## What Happened

S03 ran the live legal retrieval quality gate against TEI default and tagged ONNX using the user-provided 44-ФЗ corpus. The final artifact is FAIL. Ranking parity is mostly strong, but cross-backend cosine has severe outliers: document minimum 0.369 and query minimum 0.656 versus threshold 0.999. Worst cases are longer legal clauses/subclauses, consistent with a likely long-text/truncation divergence from the ONNX max sequence length 128 path. Tagged ONNX was cleaned up after evaluation.

## Verification

Runtime, evaluator, artifact hygiene, cleanup, and GitNexus checks completed.

## Requirements Advanced

- russian-legal-quality — Ran first live legal quality gate and identified divergence class.

## Requirements Validated

None.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

- onnx-production-quality — Tagged ONNX does not currently pass strict legal-corpus vector equivalence on the 44-ФЗ sample.

## Operational Readiness

None.

## Deviations

Live gate initially exposed evaluator ID fallback issue for unnumbered subclauses. The evaluator was fixed and the gate rerun. The final result remains FAIL due cross-backend cosine outliers.

## Known Limitations

The gate uses synthetic title/self known-item queries, not explicit human qrels. However, cross-backend cosine outliers are enough to block equivalence claims on this corpus.

## Follow-ups

S04 should record a blocking quality decision. Next technical investigation should focus on long-text/truncation behavior: ONNX model exported with max sequence length 128 diverges sharply from TEI on some longer legal clauses. Packaging/tuning should be deprioritized until this is understood.

## Files Created/Modified

- `tools/evaluate_legal_retrieval.py` — Updated evaluator ID fallback and worst cross-backend cosine diagnostics.
- `benchmark-results/fd-legal-retrieval-m015-s03.txt` — Live legal retrieval gate artifact with FAIL verdict.
