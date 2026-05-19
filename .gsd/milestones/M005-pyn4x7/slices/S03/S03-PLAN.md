# S03: Verification and commit

**Goal:** Run final verification and close the production hardening documentation milestone with a local commit.
**Demo:** After this, docs and runtime are verified and committed locally.

## Must-Haves

- `docker compose config` passes.
- `cd api && go test ./... -short` passes.
- README command snippets and hardening snippets are checked.
- GitNexus change detection is low/expected.
- GSD milestone is validated/completed and local commit created.

## Proof Level

- This slice proves: full verification gates for docs/config

## Integration Closure

Docs and GSD decision artifacts are committed with verified state.

## Verification

- Records final verification and remaining follow-ups for production hardening docs.

## Tasks

- [x] **T01: Run final documentation verification** `est:small`
  Run final verification commands for docs/config: compose config, Go tests, README snippet checks, and GitNexus change detection.
  - Files: `README.md`, `.gsd/DECISIONS.md`
  - Verify: All verification commands pass.

- [x] **T02: Validate and complete milestone** `est:small`
  Validate and complete M005 with success criteria, slice audit, and summary.
  - Files: `.gsd/milestones/M005-pyn4x7`
  - Verify: GSD validate and complete tools succeed.

- [x] **T03: Commit local hardening docs** `est:small`
  Checkpoint GSD DB, stage README/GSD artifacts, create local commit, and report final status. Do not push.
  - Files: `README.md`, `.gsd/DECISIONS.md`, `.gsd/gsd.db`, `.gsd/milestones/M005-pyn4x7`
  - Verify: git status clean except ahead marker; commit hash reported.

## Files Likely Touched

- README.md
- .gsd/DECISIONS.md
- .gsd/milestones/M005-pyn4x7
- .gsd/gsd.db
