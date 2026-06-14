---
id: T01
parent: S01
milestone: M043-dpr0cq
key_files:
  - .golangci.yml (extended with 5 Tier 1 linters)
  - benchmark-results/m043-tier1-baseline.txt (raw output)
key_decisions:
  - (none)
duration: 
verification_result: untested
completed_at: 2026-06-14T03:51:44.742Z
blocker_discovered: false
---

# T01: Baseline noise floor captured: 11 issues (1 errorlint, 3 goconst, 5 gosec, 2 unused) перед fixes

**Baseline noise floor captured: 11 issues (1 errorlint, 3 goconst, 5 gosec, 2 unused) перед fixes**

## What Happened

Добавил 5 Tier 1 linters (gosec, bodyclose, prealloc, errorlint, revive) в .golangci.yml с per-linter settings: gosec.exclusions для G107/G304, errorlint checks:all, revive 19 selected rules (exported excluded — deferred to S02), goconst path-pattern для test files (не сработал в golangci-lint v2.12.2). Запустил `golangci-lint run --config .golangci.yml ./...` 8+ раз (попытки fixing несуществующих rule names — error-returning, optimize-operands, modifies-param, modifies-val-receiver). Baseline output: 11 issues — 1 errorlint (redis.go `err == redis.Nil`), 3 goconst (test files), 5 gosec (G115 int→uint16 overflow, G304 file inclusion x3, G112 slowloris), 2 unused (`errorKey` const, `teiSubBatchSize` const). Сохранил raw output в benchmark-results/m043-tier1-baseline.txt.

## Verification

golangci-lint run --config .golangci.yml ./... captures 11 issues в baseline (warn mode). Per-linter counts recorded в Phase 1 report (после T02/T03 fixes). Raw output saved.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `.golangci.yml (extended with 5 Tier 1 linters)`
- `benchmark-results/m043-tier1-baseline.txt (raw output)`
