---
id: T03
parent: S02
milestone: M045-d0e5xq
key_files:
  - docker-compose.yaml
  - docs/same-host-embedding-service-contract.md
  - documents/tei-startup-mitigation-m045.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T11:37:43.753Z
blocker_discovered: false
---

# T03: Prepared compose/docs candidate for TEI offline cache startup without restarting runtime.

**Prepared compose/docs candidate for TEI offline cache startup without restarting runtime.**

## What Happened

Updated `docker-compose.yaml` TEI service to include `HF_HUB_OFFLINE=1` and explicit `HUGGINGFACE_HUB_CACHE=/data`. Updated the same-host service contract to document offline cache mode as startup mitigation only, not fd fallback behavior. Verified `docker compose config tei` shows the candidate while the running container remains unchanged.

## Verification

`/tmp/m045-compose-candidate-check.txt` verified compose candidate has `HF_HUB_OFFLINE=1`, `HF_HOME=/data`, and `HUGGINGFACE_HUB_CACHE=/data`; running container env still lacks `HF_HUB_OFFLINE=1`, proving no restart/recreate applied it yet.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose config tei plus docker inspect current env subset` | 0 | ✅ pass: compose candidate visible and running container unchanged | 120000ms |

## Deviations

None.

## Known Issues

The candidate is staged in compose only; S03 must prove startup behavior by controlled restart.

## Files Created/Modified

- `docker-compose.yaml`
- `docs/same-host-embedding-service-contract.md`
- `documents/tei-startup-mitigation-m045.md`
