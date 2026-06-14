---
id: T02
parent: S03
milestone: M043-dpr0cq
key_files:
  - .github/workflows/go-quality.yml
key_decisions:
  - (none)
duration: 
verification_result: untested
completed_at: 2026-06-14T05:14:38.736Z
blocker_discovered: false
---

# T02: CI govulncheck step added after golangci-lint in go-quality workflow

**CI govulncheck step added after golangci-lint in go-quality workflow**

## What Happened

Updated .github/workflows/go-quality.yml: timeout 15→20 minutes and added `Run govulncheck` step after the golangci-lint step, with working-directory api and command `go run golang.org/x/vuln/cmd/govulncheck@latest ./...`. Kept govulncheck standalone rather than as a golangci-lint plugin. Updated lint step label/comment to reflect 18 linters (baseline + Tier 1 + Tier 2).

## Verification

PyYAML parse of .github/workflows/go-quality.yml succeeded (`workflow_yaml_ok`). Text check confirms timeout-minutes:20 and `Run govulncheck` step at workflow line 89. Local govulncheck command exits 0.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| — | No verification commands discovered | — | — | — |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `.github/workflows/go-quality.yml`
