---
id: T01
parent: S03
milestone: M045-d0e5xq
key_files:
  - benchmark-results/m045-tei-local-path-startup-proof.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T12:18:58.866Z
blocker_discovered: false
---

# T01: Recorded local-path preflight, rollback plan, and proof criteria after offline candidate failed.

**Recorded local-path preflight, rollback plan, and proof criteria after offline candidate failed.**

## What Happened

The preflight artifact records that `HF_HUB_OFFLINE=1` failed, selected the cached USER-bge-m3 local snapshot path, documented rollback to Hub ID command, captured current container state, compose candidate, and local directory check before applying the restart proof.

## Verification

`benchmark-results/m045-tei-local-path-startup-proof.md` contains preflight state, local snapshot path, rollback instructions, compose candidate, and local snapshot file listing before proof execution.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `write preflight artifact and verify local snapshot directory exists inside fd_tei` | 0 | ✅ pass: preflight and rollback plan captured | 120000ms |

## Deviations

Replaced original offline-env proof path because the offline candidate timed out and stayed unhealthy.

## Known Issues

None for T01; proof execution handled by later tasks.

## Files Created/Modified

- `benchmark-results/m045-tei-local-path-startup-proof.md`
