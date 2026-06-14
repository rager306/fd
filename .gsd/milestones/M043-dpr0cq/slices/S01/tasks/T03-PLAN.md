---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T03: Tier 1 linters moved to fail mode (12 linters, 0 issues, CI timeout 15min) + Phase 1 report saved

.golangci.yml: move Tier 1 linters from warn to fail mode (severity: error default). Verify `golangci-lint run` exit 0. .github/workflows/go-quality.yml: bump timeout to 15 min (12+ linters медленнее). Проверить CI step не использует cached results иначе новые linter issues могут быть masked. docs/static-analysis-phase1-report-m043.md: baseline noise count, fix list, exclusions rationale, false positive rate per linter, before/after .golangci.yml diff. Также: confirm M041 acceptance tests (45 cases) все ещё pass.

## Inputs

- None specified.

## Expected Output

- `.golangci.yml (fail mode)`
- `.github/workflows/go-quality.yml (timeout 15min)`
- `docs/static-analysis-phase1-report-m043.md`

## Verification

golangci-lint run exit 0. CI workflow YAML valid (parseable). docs/static-analysis-phase1-report-m043.md существует, ≥2KB. go test ./api/... -short exit 0 (all M041 tests pass).
