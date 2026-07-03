---
id: T03
parent: S01
milestone: M043-dpr0cq
key_files:
  - .golangci.yml (fail mode, severity override removed)
  - .github/workflows/go-quality.yml (timeout 15min, step name + description)
  - docs/static-analysis-phase1-report-m043.md (10KB report)
key_decisions:
  - (none)
duration: 
verification_result: untested
completed_at: 2026-06-14T03:51:44.748Z
blocker_discovered: false
---

# T03: Tier 1 linters moved to fail mode (12 linters, 0 issues, CI timeout 15min) + Phase 1 report saved

**Tier 1 linters moved to fail mode (12 linters, 0 issues, CI timeout 15min) + Phase 1 report saved**

## What Happened

Move Tier 1 linters to fail mode: removed per-linter `severity: warning` override from revive config block. .github/workflows/go-quality.yml: timeout-minutes 10 → 15 (12 linters vs 7 baseline; headroom for future). Step name updated "Run GolangCI-Lint with Staticcheck + Tier 1 linters" + description. Final verification: `golangci-lint run --config .golangci.yml ./...` exit 0, 0 issues. docs/static-analysis-phase1-report-m043.md saved (10KB) с before/after metrics, fix list (16 changes), exclusions rationale, out-of-phase items (revive:exported deferred to S02, govulncheck to S03, Tier 2 to S02, Tier 3 opt-in). Files changed: .golangci.yml, .github/workflows/go-quality.yml, 11 Go source files, 1 new report file, 1 raw output artifact.

## Verification

golangci-lint run --config .golangci.yml ./... exit 0, 0 issues, fail mode (severity: error default). go test ./api/... все 5 packages pass. .github/workflows/go-quality.yml YAML valid. docs/static-analysis-phase1-report-m043.md (10068 bytes) с before/after table, fix list, exclusions rationale.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `.golangci.yml (fail mode, severity override removed)`
- `.github/workflows/go-quality.yml (timeout 15min, step name + description)`
- `docs/static-analysis-phase1-report-m043.md (10KB report)`
