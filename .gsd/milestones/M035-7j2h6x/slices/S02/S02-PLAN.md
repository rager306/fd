# S02: Hosting contract closure

**Goal:** Close M035 with outcome, decision, final guardrails, and local commit.
**Demo:** After this, M035 is verified, documented, committed locally, and ready for the next gate.

## Must-Haves

- Outcome and decision recorded.
- Final checks pass.
- Working tree clean after commit/reindex.
- No external state changes.

## Proof Level

- This slice proves: Docs/outcome leak checks, actionlint, provisioning/verifier checks, GitNexus detect, commit/reindex.

## Integration Closure

Keeps source contract, provisioning docs, and workflow input docs aligned.

## Verification

- Records the final blocker state for future sessions.

## Tasks

- [x] **T01: Record exact binary hosting decision** `est:small`
  Record GSD decision for exact ONNX binary hosting contract and update outcome if needed to reference decision.
  - Files: `.gsd/DECISIONS.md`, `benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt`
  - Verify: Decision and outcome checks pass.

- [x] **T02: Run final guardrails** `est:medium`
  Run final milestone guardrails: JSON/doc checks, py_compile/provisioning/verifiers, actionlint, Go tests/lint, tagged tests, Docker default build, leak checks, binary hygiene, port/background checks, GitNexus detect.
  - Verify: All final checks pass.

- [x] **T03: Close milestone and commit** `est:small`
  Validate and complete M035, checkpoint GSD DB, commit locally, reindex GitNexus, verify clean state, and report.
  - Verify: GSD completion, commit, reindex, post-reindex detect, clean working tree.

## Files Likely Touched

- .gsd/DECISIONS.md
- benchmark-results/fd-onnx-exact-binary-hosting-contract-m035-s01.txt
