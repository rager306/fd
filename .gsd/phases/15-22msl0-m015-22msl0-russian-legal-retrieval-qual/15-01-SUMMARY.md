---
id: S01
parent: M015-22msl0
milestone: M015-22msl0
provides:
  - Corpus hash/profile and evaluator metric contract.
requires:
  []
affects:
  - S02
  - S03
key_files:
  - tests/44-FZ-2026-articles.jsonl
  - benchmark-results/fd-legal-corpus-profile-m015-s01.txt
key_decisions:
  - Use parity/known-item retrieval metrics because no explicit qrels exist.
  - Do not claim absolute Russian/legal relevance quality from this corpus alone.
patterns_established:
  - Legal quality artifacts must not dump raw corpus text.
  - Unlabeled corpora support parity gates, not absolute relevance claims.
observability_surfaces:
  - Corpus SHA256 and schema/count profile.
  - Metric contract in task summary.
drill_down_paths:
  - .gsd/milestones/M015-22msl0/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M015-22msl0/slices/S01/tasks/T02-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-05-20T04:49:58.955Z
blocker_discovered: false
---

# S01: Corpus profile and gate design

**S01 turned the provided 44-ФЗ JSONL file into a measurable, caveated ONNX quality gate contract.**

## What Happened

S01 profiled the new 44-ФЗ JSONL corpus and defined the quality gate. The corpus has enough structured legal text for TEI-vs-ONNX retrieval parity and synthetic known-item checks, but not enough labels for absolute relevance claims. The artifact records corpus hash, counts, invalid flags, and length distributions without raw text.

## Verification

Corpus profiling and gate contract completed.

## Requirements Advanced

- russian-legal-quality — Established Russian/legal quality gate input and parity metrics.

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

No explicit query/qrel labels. Long legal parts may exceed runtime max sequence length; evaluator must record truncation/length stats.

## Follow-ups

S02 should implement the evaluator with no raw text in artifacts. It should handle long text truncation caveats and stable IDs. S03 should run TEI and tagged ONNX with isolated Redis namespaces.

## Files Created/Modified

- `benchmark-results/fd-legal-corpus-profile-m015-s01.txt` — Sanitized corpus profile for the user-provided 44-ФЗ JSONL file.
