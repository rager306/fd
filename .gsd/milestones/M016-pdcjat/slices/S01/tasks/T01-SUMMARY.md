---
id: T01
parent: S01
milestone: M016-pdcjat
key_files:
  - benchmark-results/fd-legal-retrieval-m015-s03.txt
  - tests/44-FZ-2026-articles.jsonl
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-20T05:32:48.607Z
blocker_discovered: false
---

# T01: Resolved and verified the M015 worst divergence IDs against the corpus.

**Resolved and verified the M015 worst divergence IDs against the corpus.**

## What Happened

Parsed the M015 live legal gate artifact, extracted worst cross-backend cosine document/query IDs, resolved each ID back to the 44-ФЗ JSONL corpus using the evaluator's fallback ID scheme, and verified all text SHA256 hashes match the artifact. No raw legal text was printed; only IDs, counts, chars, and hashes were used.

## Verification

A Python resolver confirmed 17 unique worst-case IDs resolve and all hashes match.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python resolve worst IDs from benchmark-results/fd-legal-retrieval-m015-s03.txt against tests/44-FZ-2026-articles.jsonl` | 0 | ✅ pass — resolved_count=17; all_hashes_match=true; worst_min_cosine=0.36948916 | 0ms |

## Deviations

None.

## Known Issues

The worst-case target set includes 17 unique document/query IDs after de-duplicating self-query IDs back to document IDs. The minimum recorded cosine is 0.36948916.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m015-s03.txt`
- `tests/44-FZ-2026-articles.jsonl`
