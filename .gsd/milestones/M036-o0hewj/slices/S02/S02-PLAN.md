# S02: Reproducible export closure

**Goal:** Record decision, run guardrails, and close M036.
**Demo:** After this, M036 is verified, summarized, committed locally, and reindexed.

## Must-Haves

- D034 or next decision recorded.
- Final checks pass.
- No external action occurred.
- Working tree clean after commit/reindex.

## Proof Level

- This slice proves: Final guardrails and GitNexus detect/reindex.

## Integration Closure

Keeps M036 docs aligned with M032/M035 source blockers.

## Verification

- Leaves an auditable closure and next-gate recommendation.

## Tasks

- [x] **T01: Record reproducible export decision** `est:small`
  Record GSD decision for planned reproducible-export workflow contract and update outcome to reference the decision.
  - Files: `.gsd/DECISIONS.md`, `benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt`
  - Verify: Decision/outcome checks pass.

- [x] **T02: Run final guardrails** `est:medium`
  Run final milestone guardrails: JSON/doc checks, py_compile/provisioning/verifiers, actionlint, Go tests/lint, tagged tests, Docker default build, leak checks, binary hygiene, port/background checks, GitNexus detect.
  - Verify: All final checks pass.

- [x] **T03: Prepare post-slice closure** `est:small`
  Record closure-ordering correction and defer milestone validation/completion/checkpoint/commit/reindex to the post-slice sequence.
  - Verify: Task records that post-slice closure will run after S02 completion.

## Files Likely Touched

- .gsd/DECISIONS.md
- benchmark-results/fd-onnx-reproducible-export-contract-m036-s01.txt
