# S03: Add GolangCI-Lint and Staticcheck gate

**Goal:** Add GolangCI-Lint configuration using Staticcheck and common Go analyzers, then fix any reported issues.
**Demo:** After this, `golangci-lint run` is configured and passes with Staticcheck enabled.

## Must-Haves

- `.golangci.yml` or equivalent config is added.
- Config includes Staticcheck through GolangCI-Lint.
- `golangci-lint run` passes for the Go module.
- Any code changes from lint fixes preserve tests.

## Proof Level

- This slice proves: configured lint command evidence

## Integration Closure

Creates a repeatable static-analysis gate for future Go changes.

## Verification

- Lint failures become standardized quality signals in project docs and CI-ready config.

## Tasks

- [x] **T01: Add GolangCI-Lint config** `est:small`
  Add root GolangCI-Lint config with Staticcheck and common analyzers suitable for the existing Go module under api/.
  - Files: `.golangci.yml`
  - Verify: Config file exists and references staticcheck.

- [x] **T02: Run configured lint gate** `est:medium`
  Run GolangCI-Lint via reproducible go-run/install path against api module, then inspect reported issues.
  - Files: `.golangci.yml`, `api`
  - Verify: Lint command runs; failures recorded or pass.

- [x] **T03: Fix lint findings and verify** `est:medium`
  Fix any lint issues reported by GolangCI-Lint/Staticcheck and rerun tests/lint.
  - Files: `api`
  - Verify: `go test` and configured lint command pass.

## Files Likely Touched

- .golangci.yml
- api
