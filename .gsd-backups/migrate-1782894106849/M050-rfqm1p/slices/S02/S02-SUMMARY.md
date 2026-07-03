---
id: S02
parent: M050-rfqm1p
milestone: M050-rfqm1p
provides:
  - Validated R044 real-service e2e proof.
  - Reusable `tests/integration` suite for future agents.
requires:
  []
affects:
  - S03
  - S04
key_files:
  - tests/integration/api_test.go
  - benchmark-results/m050-s02-docker-e2e.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Use temporary untracked Compose override for local authenticated proof instead of reading or overwriting `api/.env`.
  - Treat `/metrics` as authenticated diagnostics in integration tests.
patterns_established:
  - Two-mode e2e: no-key public/fail-closed mode and authenticated full runtime mode.
  - Explicit integration secret variable `FD_INTEGRATION_API_KEY`.
observability_surfaces:
  - `benchmark-results/m050-s02-docker-e2e.md` records e2e contract and pass summary.
drill_down_paths:
  - .gsd/milestones/M050-rfqm1p/slices/S02/tasks/T01-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S02/tasks/T02-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S02/tasks/T03-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S02/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-06-15T14:49:27.309Z
blocker_discovered: false
---

# S02: Docker e2e suite for current service

**Поддерживаемый Docker Compose e2e suite создан и прошёл authenticated proof против реальных `fd_api`, `fd_redis`, `fd_tei`.**

## What Happened

S02 расширил `tests/integration` из минимального actuality layer в текущий black-box e2e smoke. Suite проверяет public probes/health context, fail-closed auth для protected endpoints, authenticated metrics, embeddings dimensions and batch behavior, validation errors, missing-model compatibility, and cache HIT/flush/delete invalidation. Authenticated proof пересоздал API с temporary nonprinted Compose key and ran the suite with the matching `FD_INTEGRATION_API_KEY`; all checks passed.

## Verification

No-key mode passed with 5 checks. Authenticated Docker Compose mode passed with `SUMMARY pass=9 fail=0 skip=0`. Evidence: `benchmark-results/m050-s02-docker-e2e.md`.

## Requirements Advanced

None.

## Requirements Validated

- R044 — Authenticated Docker Compose e2e run passed with `SUMMARY pass=9 fail=0 skip=0`; evidence artifact `benchmark-results/m050-s02-docker-e2e.md`.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

The initial test assumed `/metrics` was public; runtime returned 401, so metrics was correctly moved into authenticated diagnostics.

## Known Limitations

The suite needs `FD_INTEGRATION_API_KEY` for full authenticated coverage; no-key mode intentionally skips protected happy paths.

## Follow-ups

Proceed to S03 mutation baseline for critical packages.

## Files Created/Modified

- `tests/integration/api_test.go` — Expanded root integration suite into current-service Docker e2e checks.
- `benchmark-results/m050-s02-docker-e2e.md` — Recorded contract, secret handling, no-key and authenticated results, and final verdict.
- `.gsd/REQUIREMENTS.md` — R044 marked validated.
