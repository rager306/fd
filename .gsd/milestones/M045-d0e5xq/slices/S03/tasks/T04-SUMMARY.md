---
id: T04
parent: S03
milestone: M045-d0e5xq
key_files:
  - benchmark-results/m045-tei-local-path-startup-proof.md
  - documents/tei-startup-mitigation-m045.md
  - docs/same-host-embedding-service-contract.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T12:18:58.874Z
blocker_discovered: false
---

# T04: Validated R028 and updated docs to local snapshot startup posture.

**Validated R028 and updated docs to local snapshot startup posture.**

## What Happened

Updated mitigation and contract docs to record local snapshot path as the validated startup posture and rejected `HF_HUB_OFFLINE=1`. Updated R028 to validated with proof evidence. Final compose config shows local path and no `HF_HUB_OFFLINE`; final fd smoke passed.

## Verification

R028 validated via gsd_requirement_update. `docker compose config tei` contains the local snapshot path and no `HF_HUB_OFFLINE`; final fd smoke passed; `go test ./...` passed with 270 tests.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `docker compose config tei local-path assertion; fd final smoke; cd api && go test ./...` | 0 | ✅ pass: final config/runtime/go test checks passed | 120000ms |

## Deviations

None.

## Known Issues

If the cached snapshot revision changes, compose must be updated to the new local path and proof rerun.

## Files Created/Modified

- `benchmark-results/m045-tei-local-path-startup-proof.md`
- `documents/tei-startup-mitigation-m045.md`
- `docs/same-host-embedding-service-contract.md`
- `.gsd/REQUIREMENTS.md`
