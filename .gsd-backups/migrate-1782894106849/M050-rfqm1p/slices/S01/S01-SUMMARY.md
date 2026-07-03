---
id: S01
parent: M050-rfqm1p
milestone: M050-rfqm1p
provides:
  - Current existing-test baseline for S02 Docker e2e work.
  - Validated R043 evidence.
requires:
  []
affects:
  - S02
  - S03
  - S04
key_files:
  - tests/integration/api_test.go
  - tests/integration/go.mod
  - benchmark-results/m050-s01-test-actuality.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Use `FD_INTEGRATION_API_KEY` instead of `FD_API_KEY` for root integration tests to avoid accidental stale shell secret coupling.
patterns_established:
  - Existing test actuality audit before adding new test layers.
  - Explicit integration-test secret variable separate from service runtime env var.
observability_surfaces:
  - `benchmark-results/m050-s01-test-actuality.md` records stale/fixed/deferred classifications and command evidence.
drill_down_paths:
  - .gsd/milestones/M050-rfqm1p/slices/S01/tasks/T01-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S01/tasks/T02-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S01/tasks/T03-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S01/tasks/T04-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-06-15T14:41:02.345Z
blocker_discovered: false
---

# S01: Existing test actuality audit

**Существующие тесты проверены на актуальность; stale root integration layer исправлен и основной `api` baseline остался зелёным.**

## What Happened

S01 провёл инвентарь существующих тестов и verification scripts перед добавлением новых e2e/mutation слоёв. Регулярные `api` проверки оказались актуальными и прошли. Отдельный `tests/integration` слой был stale: он не был Go module, предполагал protected embeddings без auth и мог использовать случайный shell `FD_API_KEY`. Срез добавил standalone module, привёл integration test к текущему fail-closed auth контракту и сделал protected checks зависимыми от явного `FD_INTEGRATION_API_KEY`.

## Verification

Final commands passed: `cd api && go test ./...` returned 295 passed in 10 packages; `cd tests/integration && go test -v .` returned 2 passed in 1 package. Artifact: `benchmark-results/m050-s01-test-actuality.md`.

## Requirements Advanced

None.

## Requirements Validated

- R043 — `benchmark-results/m050-s01-test-actuality.md`, final `cd api && go test ./...`, final `cd tests/integration && go test -v .`.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

Full authenticated Docker Compose e2e coverage is intentionally deferred to S02; S01 only corrected the existing root integration layer.

## Known Limitations

`tests/integration` protected happy-path checks skip without `FD_INTEGRATION_API_KEY`. This avoids reading or printing local secrets; S02 will introduce a maintained authenticated e2e flow.

## Follow-ups

Proceed to S02 to build the current-service Docker Compose e2e suite required by R044.

## Files Created/Modified

- `tests/integration/api_test.go` — Updated root integration tests for current auth posture, base URL env, explicit integration key, and fail-closed unauthenticated embeddings check.
- `tests/integration/go.mod` — Added standalone integration test module using Go 1.25.0.
- `benchmark-results/m050-s01-test-actuality.md` — Recorded inventory, stale findings, fixes, commands, and S02 deferrals.
- `.gsd/REQUIREMENTS.md` — R043 marked validated with evidence.
