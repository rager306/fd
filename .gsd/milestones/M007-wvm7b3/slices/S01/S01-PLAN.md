# S01: CI workflow design

**Goal:** Research current GitHub Actions syntax and inspect repo workflow support before writing files.
**Demo:** After this, workflow syntax/action choices are grounded in current docs and repo constraints.

## Must-Haves

- GitHub Actions workflow syntax docs are consulted.
- Existing workflow/scripts availability is inspected.
- CI design records triggers, permissions, Go setup, cache, and lint command.
- No external GitHub state changes are made.

## Proof Level

- This slice proves: docs plus repo inspection

## Integration Closure

Workflow implementation will be based on current docs and actual repo constraints.

## Verification

- Records CI design and local/remote validation limits.

## Tasks

- [x] **T01: Consult GitHub Actions syntax docs** `est:small`
  Fetch current GitHub Actions workflow syntax docs for `name`, `on`, `permissions`, and `jobs` basics.
  - Verify: Relevant docs snippets summarized.

- [x] **T02: Inspect local CI support** `est:small`
  Inspect repository for existing `.github/workflows`, helper scripts such as `scripts/ci_monitor.cjs`, and README quality command parity.
  - Files: `README.md`, `.golangci.yml`
  - Verify: Existing workflow/support status recorded.

- [x] **T03: Record workflow design** `est:small`
  Record minimal workflow design: triggers, permissions, Go setup, cache/dependency strategy, and exact test/lint commands.
  - Files: `README.md`, `.golangci.yml`
  - Verify: Design summary recorded.

## Files Likely Touched

- README.md
- .golangci.yml
