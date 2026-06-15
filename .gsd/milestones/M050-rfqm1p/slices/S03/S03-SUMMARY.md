---
id: S03
parent: M050-rfqm1p
milestone: M050-rfqm1p
provides:
  - Validated R045 mutation baseline.
  - Exact local command for future mutation checks.
requires:
  []
affects:
  - S04
key_files:
  - benchmark-results/m050-s03-mutation-baseline.md
  - .gsd/REQUIREMENTS.md
key_decisions:
  - Use `avito-tech/go-mutesting` for current bounded baseline.
  - Do not enable mutation as mandatory CI hard gate in M050.
patterns_established:
  - Bounded critical-file mutation baseline with explicit score and scope.
observability_surfaces:
  - `benchmark-results/m050-s03-mutation-baseline.md` records runner, command, score, and policy.
drill_down_paths:
  - .gsd/milestones/M050-rfqm1p/slices/S03/tasks/T01-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S03/tasks/T02-SUMMARY.md
  - .gsd/milestones/M050-rfqm1p/slices/S03/tasks/T03-SUMMARY.md
duration: ""
verification_result: passed
completed_at: 2026-06-15T14:55:16.221Z
blocker_discovered: false
---

# S03: Mutation baseline for critical packages

**Рабочий bounded mutation baseline создан: score 1.0 на выбранных cache, handlers и lifecycle файлах.**

## What Happened

S03 проверил Go mutation tooling, выбрал свежий `avito-tech/go-mutesting` runner, подтвердил smoke на `cache/hash.go`, затем запустил bounded critical baseline на `cache/hash.go`, `cache/keys.go`, `handlers/cache.go`, `handlers/health.go`, and `lifecycle/state.go`. Все мутанты в scope были убиты; survivors не обнаружены. Baseline зафиксирован как local/manual informational gate, не как обязательный CI hard gate.

## Verification

Mutation critical baseline exited 0 with `The mutation score is 1.000000 (143 passed, 0 failed, 4 duplicated, 0 skipped, total is 143)`. Fresh `cd api && go test ./...` passed with 295 tests in 10 packages.

## Requirements Advanced

None.

## Requirements Validated

- R045 — Mutation score 1.000000 on selected critical files; artifact `benchmark-results/m050-s03-mutation-baseline.md`.

## New Requirements Surfaced

None.

## Requirements Invalidated or Re-scoped

None.

## Operational Readiness

None.

## Deviations

No repo-wide mutation run was attempted; this was intentionally bounded for runtime and signal quality.

## Known Limitations

The selected runner requires Go >= 1.25.5 and auto-switches to go1.25.11. CI hard gating is deferred pending runner pinning/toolchain caching.

## Follow-ups

S04 should document mutation as local/manual informational gate and list the exact command.

## Files Created/Modified

- `benchmark-results/m050-s03-mutation-baseline.md` — Recorded mutation tooling choice, smoke, bounded baseline, score and policy.
- `.gsd/REQUIREMENTS.md` — R045 marked validated.
