---
id: M050-rfqm1p
title: "Current test truth and stronger gates"
status: complete
completed_at: 2026-06-15T14:58:50.097Z
key_decisions:
  - Use `FD_INTEGRATION_API_KEY` for integration client auth rather than accidental shell `FD_API_KEY`.
  - Treat `/metrics` as authenticated diagnostics after runtime returned 401 without auth.
  - Use temporary untracked Compose override for authenticated e2e proof instead of reading or overwriting `api/.env`.
  - Use `avito-tech/go-mutesting` as the bounded mutation runner and keep mutation local/manual for now.
  - Do not add Docker e2e or mutation to CI in M050; document them as manual/local gates.
key_files:
  - tests/integration/api_test.go
  - tests/integration/go.mod
  - README.md
  - benchmark-results/m050-s01-test-actuality.md
  - benchmark-results/m050-s02-docker-e2e.md
  - benchmark-results/m050-s03-mutation-baseline.md
  - benchmark-results/m050-s04-test-gates-closure.md
  - .gsd/milestones/M050-rfqm1p/M050-rfqm1p-VALIDATION.md
lessons_learned:
  - Existing tests can be stale even when main package tests are green; root integration was outside the module and assumed pre-auth behavior.
  - No-key and authenticated integration modes should be explicit to avoid leaking or accidentally using stale secrets.
  - Mutation score is useful only with a bounded, reproducible scope and clear runner/toolchain policy.
---

# M050-rfqm1p: Current test truth and stronger gates

**M050 сделал тестовую систему fd актуальной и сильнее: stale тесты исправлены, Docker e2e suite создан, mutation baseline измерен, команды документированы.**

## What Happened

Milestone M050 began with the user-directed constraint that all existing tests must first be checked for actuality against the latest service version. S01 inventoried the existing suite, found that the root integration layer was stale and not runnable, corrected it as a standalone module, and validated R043. S02 expanded `tests/integration` into a real auth-aware Docker Compose e2e suite and verified it against `fd_api`, `fd_redis`, and `fd_tei` with a temporary nonprinted key, validating R044. S03 selected a current Go mutation runner and established a bounded critical-file mutation baseline with score 1.0, validating R045. S04 updated README with exact commands and gate policy, then ran final verification across tests, lint, govulncheck, and integration no-key mode.

## Success Criteria Results

- Existing tests classified and stale root integration fixed: PASS, `benchmark-results/m050-s01-test-actuality.md`.
- API suite current and passing: PASS, final `cd api && go test ./...` 295 passed.
- Docker e2e suite current and authenticated: PASS, `benchmark-results/m050-s02-docker-e2e.md`, pass=9 fail=0 skip=0.
- Mutation baseline exists: PASS, `benchmark-results/m050-s03-mutation-baseline.md`, score 1.000000 over 143 mutants in scope.
- Commands documented: PASS, README Development section and `benchmark-results/m050-s04-test-gates-closure.md`.

## Definition of Done Results

- No stale existing test layer remains silently broken: met.
- Commands are reproducible with documented prerequisites: met.
- Evidence artifacts saved under `benchmark-results/`: met.
- Requirements R043-R045 validated: met.
- Worktree cleanliness still pending final local git check after GSD writes; no known untracked runtime secrets are part of milestone output.

## Requirement Outcomes

- R043 validated by S01 actuality audit and final Go/root integration runs.
- R044 validated by S02 authenticated Docker Compose e2e proof.
- R045 validated by S03 bounded mutation baseline.

## Deviations

M050 did not wire heavy Docker e2e or mutation into CI; this was intentional to avoid secret/runtime cost surprises.

## Follow-ups

Optional future work: add a manual GitHub Actions workflow for Docker e2e after CI secret/runtime policy is chosen; pin mutation runner version if making it a recurring gate.
