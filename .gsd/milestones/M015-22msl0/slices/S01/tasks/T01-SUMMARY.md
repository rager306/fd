---
id: T01
parent: S01
milestone: M015-22msl0
key_files:
  - tests/44-FZ-2026-articles.jsonl
  - benchmark-results/fd-legal-corpus-profile-m015-s01.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-20T04:49:29.061Z
blocker_discovered: false
---

# T01: Profiled the 44-ФЗ JSONL legal corpus and wrote a sanitized profile artifact.

**Profiled the 44-ФЗ JSONL legal corpus and wrote a sanitized profile artifact.**

## What Happened

Profiled the user-provided 44-ФЗ JSONL corpus without dumping raw legal text. The file has 94 article records, 668 parts, 912 clauses, and 272 subclauses. Some entries are marked invalid. Part text length is highly skewed, with p95 3388 chars and max 42874 chars, so truncation caveats must be recorded in the evaluator. The corpus SHA256 is `de03cda6b266085a9b1f2376afcb9dffbb00fec922dee1f1553cadcfb6d03869`.

## Verification

Profile artifact exists, includes counts/hash/schema/length stats, and explicitly excludes raw legal text.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python profile tests/44-FZ-2026-articles.jsonl > benchmark-results/fd-legal-corpus-profile-m015-s01.txt` | 0 | ✅ pass — profile artifact written | 0ms |
| 2 | `read benchmark-results/fd-legal-corpus-profile-m015-s01.txt` | 0 | ✅ pass — counts/hash/no raw text caveat present | 0ms |

## Deviations

None.

## Known Issues

The corpus has no explicit qrels; it can support parity/known-item checks, not absolute human relevance scoring.

## Files Created/Modified

- `tests/44-FZ-2026-articles.jsonl`
- `benchmark-results/fd-legal-corpus-profile-m015-s01.txt`
