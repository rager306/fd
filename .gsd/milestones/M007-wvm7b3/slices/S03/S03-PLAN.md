# S03: Document and verify CI workflow

**Goal:** Document CI quality gate, run final local verification, complete GSD milestone, and commit locally.
**Demo:** After this, CI docs are updated and the milestone is locally committed.

## Must-Haves

- README mentions CI quality gate and local command parity.
- Workflow YAML parses or passes available local checks.
- `cd api && go test ./... -short` passes.
- Pinned GolangCI-Lint command passes.
- GitNexus change detection is expected.
- GSD milestone complete and local commit created.

## Proof Level

- This slice proves: local final gates

## Integration Closure

Leaves repository ready to push and observe remote CI once the user explicitly approves.

## Verification

- README and GSD summaries explain local vs remote CI verification.

## Tasks

- [x] **T01: Document CI workflow** `est:small`
  Update README quality tooling section to mention `.github/workflows/go-quality.yml`, triggers, and remote run pending push.
  - Files: `README.md`
  - Verify: README contains workflow path and CI trigger notes.

- [x] **T02: Run final CI workflow verification** `est:medium`
  Run final local verification: workflow parse, README snippets, Go tests, GolangCI-Lint, and GitNexus change detection.
  - Files: `.github/workflows/go-quality.yml`, `README.md`, `api`
  - Verify: All local checks pass.

- [x] **T03: Close and commit CI milestone** `est:small`
  Validate/complete M007, checkpoint GSD DB, stage workflow/docs/GSD artifacts, create local commit, and report status. Do not push.
  - Files: `.github/workflows/go-quality.yml`, `README.md`, `.gsd/gsd.db`, `.gsd/milestones/M007-wvm7b3`
  - Verify: git status clean except ahead marker; commit hash reported.

## Files Likely Touched

- README.md
- .github/workflows/go-quality.yml
- api
- .gsd/gsd.db
- .gsd/milestones/M007-wvm7b3
