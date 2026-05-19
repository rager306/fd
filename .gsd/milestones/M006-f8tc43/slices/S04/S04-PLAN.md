# S04: Document and verify quality tooling

**Goal:** Document test/lint usage, run final gates, complete GSD milestone, and create local commit.
**Demo:** After this, README documents quality commands and the milestone is committed locally.

## Must-Haves

- README documents `go test`, Testify, and GolangCI-Lint usage.
- Final Go tests pass.
- Final GolangCI-Lint run passes.
- GitNexus change detection is low/expected or risk is documented.
- GSD milestone complete and local commit created.

## Proof Level

- This slice proves: full verification gates

## Integration Closure

Leaves project with committed, documented quality tooling.

## Verification

- Quality commands are visible to future agents in README and GSD summaries.

## Tasks

- [x] **T01: Document quality tooling commands** `est:small`
  Update README development section with Go test command, Testify usage pattern, and pinned GolangCI-Lint command with Staticcheck enabled by config.
  - Files: `README.md`
  - Verify: README contains test/lint/Testify/Staticcheck snippets.

- [x] **T02: Run final quality verification** `est:medium`
  Run final verification: Go tests, pinned GolangCI-Lint command, README snippet checks, and GitNexus change detection.
  - Files: `README.md`, `.golangci.yml`, `api`
  - Verify: All commands pass.

- [x] **T03: Close and commit quality tooling milestone** `est:small`
  Validate and complete M006, checkpoint GSD DB, stage relevant files, create local commit, and report status. Do not push.
  - Files: `.gsd/gsd.db`, `.gsd/milestones/M006-f8tc43`, `README.md`, `.golangci.yml`, `api`
  - Verify: git status clean except ahead marker; commit hash reported.

## Files Likely Touched

- README.md
- .golangci.yml
- api
- .gsd/gsd.db
- .gsd/milestones/M006-f8tc43
