---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T01: Baseline noise floor captured: 11 issues (1 errorlint, 3 goconst, 5 gosec, 2 unused) перед fixes

.golangci.yml: добавить Tier 1 linters (gosec, bodyclose, prealloc, errorlint, revive) с per-linter settings (gosec: G107 warn для URL-from-env, revive: explicit rules list, errorlint: check errorf formatting). Запустить `golangci-lint run --config .golangci.yml ./...` локально. Собрать output: per-linter issue count, false positive classification, fix vs exclude decision. Это baseline для Phase 1 report.

## Inputs

- None specified.

## Expected Output

- `.golangci.yml (extended)`
- `tools/lint-baseline.sh (capture output to file)`
- `benchmark-results/m043-tier1-baseline.txt (raw output)`

## Verification

golangci-lint run exit code ≠ 0 (issues captured), но без падений compile. Per-linter counts recorded.
