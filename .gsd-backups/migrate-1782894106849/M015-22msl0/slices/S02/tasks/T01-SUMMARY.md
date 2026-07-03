---
id: T01
parent: S02
milestone: M015-22msl0
key_files:
  - tools/evaluate_legal_retrieval.py
key_decisions:
  - Evaluator uses non-invalid clause/subclause documents first, falling back to parts only if needed.
  - Queries include article-title known-item queries and self-document queries; artifacts contain IDs/hashes/metrics only.
  - Default thresholds are strict parity gates: top1 agreement 0.90, mean overlap@5 0.90, ONNX recall ratio 0.98, min cross-backend cosine 0.999.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:53:13.951Z
blocker_discovered: false
---

# T01: Implemented the sanitized Russian/legal retrieval parity evaluator.

**Implemented the sanitized Russian/legal retrieval parity evaluator.**

## What Happened

Implemented `tools/evaluate_legal_retrieval.py`. The tool loads the structured JSONL corpus, builds sanitized candidate documents and queries, supports dry-run mode, embeds documents and queries through TEI and ONNX API URLs, compares cross-backend cosine, top-1 agreement, top-k overlap, TEI/ONNX recall@k, MRR, and worst query IDs. It renders markdown artifacts without raw legal text.

## Verification

`python3 -m py_compile tools/evaluate_legal_retrieval.py` passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/evaluate_legal_retrieval.py` | 0 | ✅ pass | 0ms |

## Deviations

None.

## Known Issues

The evaluator performs synthetic known-item checks, not human-labeled legal relevance. It should be extended if explicit qrels are added.

## Files Created/Modified

- `tools/evaluate_legal_retrieval.py`
