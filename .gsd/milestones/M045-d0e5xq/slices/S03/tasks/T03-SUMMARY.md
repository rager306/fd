---
id: T03
parent: S03
milestone: M045-d0e5xq
key_files:
  - benchmark-results/m045-tei-local-path-startup-proof.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T12:18:58.873Z
blocker_discovered: false
---

# T03: Verified fd and direct TEI runtime after local-path startup.

**Verified fd and direct TEI runtime after local-path startup.**

## What Happened

After TEI became healthy, ran fd health, fd ready, fd embedding smoke, and direct TEI embedding smoke. fd still reports backend `tei`, model `deepvk/USER-bge-m3`, dimensions 1024, and both embedding paths returned 1024-dimensional vectors.

## Verification

`benchmark-results/m045-tei-local-path-startup-proof.md` includes smoke results. Additional final fd smoke also passed after docs updates.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `fd /health, fd /ready, fd /v1/embeddings, direct TEI /embeddings after local path startup` | 0 | ✅ pass: fd and TEI smoke checks passed | 120000ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `benchmark-results/m045-tei-local-path-startup-proof.md`
