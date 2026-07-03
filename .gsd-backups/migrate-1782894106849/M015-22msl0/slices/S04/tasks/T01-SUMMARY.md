---
id: T01
parent: S04
milestone: M015-22msl0
key_files:
  - benchmark-results/fd-legal-retrieval-m015-summary.txt
  - benchmark-results/fd-legal-retrieval-m015-s03.txt
key_decisions:
  - The gate verdict is FAIL because cosine equivalence has severe legal-corpus outliers despite mostly strong ranking parity.
  - Next priority should be long-text legal quality divergence investigation, not ONNX packaging/tuning.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:05:39.439Z
blocker_discovered: false
---

# T01: Summarized the M015 legal quality gate as a FAIL with next-step implications.

**Summarized the M015 legal quality gate as a FAIL with next-step implications.**

## What Happened

Wrote the M015 quality gate summary. It records the 44-ФЗ corpus hash, selected doc/query counts, TEI/ONNX endpoints, key metrics, fail verdict, interpretation, and decision impact. It explicitly states that TEI remains default and that ONNX packaging/tuning should wait until long-text legal divergence is investigated.

## Verification

Summary artifact exists and includes verdict, metrics, caveats, and no raw legal text.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python summary generator from benchmark-results/fd-legal-retrieval-m015-s03.txt` | 0 | ✅ pass — summary artifact written | 0ms |
| 2 | `read benchmark-results/fd-legal-retrieval-m015-summary.txt` | 0 | ✅ pass — verdict/metrics/caveats present | 0ms |

## Deviations

None.

## Known Issues

No explicit qrels; summary is parity/known-item evidence only. However, vector equivalence failure is enough to block production-readiness.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m015-summary.txt`
- `benchmark-results/fd-legal-retrieval-m015-s03.txt`
