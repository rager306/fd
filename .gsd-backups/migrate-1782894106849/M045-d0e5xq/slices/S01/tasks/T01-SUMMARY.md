---
id: T01
parent: S01
milestone: M045-d0e5xq
key_files:
  - documents/tei-startup-recon-m045.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T11:30:47.891Z
blocker_discovered: false
---

# T01: Captured current TEI runtime state without restarting containers.

**Captured current TEI runtime state without restarting containers.**

## What Happened

Collected read-only Docker/container evidence for fd_tei, fd_api, and fd_redis; captured TEI image, digest, command, safe env subset, health state, fd health/readiness, fd embedding smoke, and direct TEI embedding smoke. Wrote the evidence to `documents/tei-startup-recon-m045.md`. No restart, recreate, docker run, or compose mutation was performed.

## Verification

`documents/tei-startup-recon-m045.md` contains current runtime state; fd `/health`, `/ready`, fd `/v1/embeddings`, and direct TEI `/embeddings` returned HTTP 200 with 1024-dimensional embeddings.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read-only Python collector using docker inspect, docker compose config/logs, fd health/ready/embedding smoke, and direct TEI embedding smoke` | 0 | ✅ pass: runtime state and smoke evidence captured | 180000ms |

## Deviations

None.

## Known Issues

Smoke proof confirms steady-state health only; startup mitigation still requires controlled proof in later slices.

## Files Created/Modified

- `documents/tei-startup-recon-m045.md`
