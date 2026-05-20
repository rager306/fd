---
id: T02
parent: S01
milestone: M015-22msl0
key_files:
  - benchmark-results/fd-legal-corpus-profile-m015-s01.txt
key_decisions:
  - Because the corpus has no explicit qrels, the gate measures TEI-vs-ONNX retrieval parity plus synthetic known-item retrieval, not absolute legal relevance.
  - Candidate documents should use stable IDs and sanitized hashes/lengths in artifacts, not raw text.
  - Initial thresholds should be strict for parity: high top-k overlap and no large ONNX degradation in synthetic known-item metrics; final threshold may be refined after dry-run distribution is known.
duration: 
verification_result: passed
completed_at: 2026-05-20T04:49:40.561Z
blocker_discovered: false
---

# T02: Defined the first-pass Russian/legal retrieval parity contract for the unlabeled 44-ФЗ corpus.

**Defined the first-pass Russian/legal retrieval parity contract for the unlabeled 44-ФЗ corpus.**

## What Happened

Defined the M015 retrieval parity contract. The evaluator should build legal document candidates from non-invalid clauses/parts where practical, derive query IDs from article titles and/or deterministic snippets, embed the same docs and queries through TEI and tagged ONNX endpoints, and compare ranking parity. Metrics should include top-1 agreement, top-3/top-5 overlap, Spearman/Kendall-like rank agreement over top candidates, cosine drift for corresponding embeddings, and synthetic known-item recall/MRR where query-to-document IDs are deterministic. Artifacts must include corpus hash, counts, runtime config, metric summaries, and worst-case IDs only, not raw legal text.

## Verification

Task summary records metrics, caveats, and artifact constraints for S02 evaluator implementation.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `GSD task summary metric contract` | 0 | ✅ pass — metrics and caveats documented | 0ms |

## Deviations

None.

## Known Issues

Synthetic known-item queries from article titles or document snippets can be biased and are weaker than human-labeled legal questions. A future dataset should add query/qrel labels.

## Files Created/Modified

- `benchmark-results/fd-legal-corpus-profile-m015-s01.txt`
