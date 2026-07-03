# S02: Add Go quality workflow

**Goal:** Add GitHub Actions workflow for Go test and GolangCI-Lint quality gate.
**Demo:** After this, `.github/workflows/go-quality.yml` runs tests and GolangCI-Lint locally-equivalent commands.

## Must-Haves

- Workflow exists under `.github/workflows/`.
- Workflow uses Go 1.25.x setup.
- Workflow runs `go test ./... -short` in `api/`.
- Workflow runs pinned GolangCI-Lint v2.12.2 with `../.golangci.yml`.
- Workflow has minimal permissions and path-aware triggers.

## Proof Level

- This slice proves: file validation plus command parity

## Integration Closure

CI mirrors M006 local quality tooling.

## Verification

- Future pushes/PRs get visible Go quality gate status.

## Tasks

- [x] **T01: Create Go quality workflow** `est:small`
  Create `.github/workflows/go-quality.yml` from the S01 design.
  - Files: `.github/workflows/go-quality.yml`
  - Verify: Workflow file exists and contains expected jobs/commands.

- [x] **T02: Verify workflow command parity** `est:small`
  Check workflow syntax structurally and verify command parity against README/local quality commands.
  - Files: `.github/workflows/go-quality.yml`, `README.md`
  - Verify: YAML parse/snippet checks pass.

## Files Likely Touched

- .github/workflows/go-quality.yml
- README.md
