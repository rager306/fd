---
id: S02
parent: M015-22msl0
milestone: M015-22msl0
provides:
  - Live evaluator for S03 TEI-vs-tagged-ONNX quality gate.
requires:
  []
affects:
  - S03
key_files:
  - tools/evaluate_legal_retrieval.py
  - benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt
key_decisions:
  - Use sanitized markdown output with full config and JSON metrics.
  - Use IDs/hashes/counts only; never raw legal text in artifacts.
patterns_established:
  - Quality-gate tooling must have dry-run mode before expensive runtime calls.
  - Artifacts use IDs and metrics only, not raw legal text.
observability_surfaces:
  - Dry-run artifact with corpus hash, selected counts, thresholds, and endpoint labels.
  - Evaluator blocked-artifact behavior for runtime failures.
drill_down_paths:
  - .gsd/milestones/M015-22msl0/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M015-22msl0/slices/S02/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T04:54:34.871Z
blocker_discovered: false
---

# S02: Legal retrieval evaluator

**S02 delivered the reusable Russian/legal retrieval parity evaluator.**

## What Happened

S02 implemented and dry-run verified `tools/evaluate_legal_retrieval.py`. The evaluator supports TEI and ONNX API URLs, corpus selection, thresholds, dry-run mode, ranking metrics, cross-backend cosine, and sanitized markdown artifacts. Dry-run proved corpus parsing and artifact hygiene before live API calls.

## Verification

Compile, dry-run, hygiene, and GitNexus checks passed.

## Requirements Advanced

- russian-legal-quality — Implemented repeatable Russian/legal retrieval parity evaluator.

## Requirements Validated

- evaluator-hygiene — Dry-run artifact includes no raw legal text and records corpus hash/config.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Dry-run query count expectation was corrected after observing selected document coverage. No evaluator code change was needed for that issue.

## Known Limitations

Evaluator quality is bounded by synthetic title/self queries because no human qrels exist. It is a parity gate, not a final relevance benchmark.

## Follow-ups

S03 should start tagged ONNX on port 18000 with `EMBEDDING_CACHE_VERSION=m015-onnx-legal-quality`, confirm TEI health, then run evaluator live. Consider starting with 256 docs and 64 self queries plus available title queries to keep runtime bounded.

## Files Created/Modified

- `tools/evaluate_legal_retrieval.py` — Sanitized TEI-vs-ONNX legal retrieval parity evaluator.
- `benchmark-results/fd-legal-retrieval-dry-run-m015-s02.txt` — Dry-run artifact proving corpus parsing and sanitized output shape.
