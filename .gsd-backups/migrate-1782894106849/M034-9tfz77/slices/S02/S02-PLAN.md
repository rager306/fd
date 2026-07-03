# S02: Hosted workflow input contract documentation

**Goal:** Document safe hosted workflow input contract and close M034 locally.
**Demo:** After this, docs/outcome/decision describe safe future dispatch inputs and remaining blockers, then M034 is verified and committed locally.

## Must-Haves

- Docs list required/optional workflow inputs and safety policy.
- Outcome records that no workflow was dispatched.
- Final verification passes; no external state changes.

## Proof Level

- This slice proves: Docs/outcome checks plus project guardrails, GitNexus detect, commit/reindex.

## Integration Closure

Keeps workflow, provisioning docs, and source-contract outcomes aligned.

## Verification

- Future operator can prepare a workflow dispatch without signed URL leaks or overclaiming readiness.

## Tasks

- [x] **T01: Document hosted workflow inputs** `est:small`
  Update provisioning docs and artifact README with safe manual workflow dispatch input contract: required ONNX/native sources, optional tokenizer/runtime sources, optional runtime sha override, no signed/plain secret URLs, and remaining exact model blocker.
  - Files: `docs/onnx-artifacts/PROVISIONING.md`, `docs/onnx-artifacts/README.md`
  - Verify: Docs contain required/optional input policy and no overclaiming.

- [x] **T02: Record workflow alignment outcome and decision** `est:small`
  Create M034 outcome artifact and decision capturing workflow input alignment and remaining blockers.
  - Files: `benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt`, `.gsd/DECISIONS.md`
  - Verify: Outcome/decision checks pass and contain no leak markers.

- [x] **T03: Verify and close milestone** `est:medium`
  Run final guardrails, complete M034, checkpoint DB, commit locally, reindex GitNexus, and report state.
  - Verify: actionlint, py_compile/provisioning/verifier/export-contract, Go checks, docs leak checks, tracked binary hygiene, GitNexus detect, commit, reindex.

## Files Likely Touched

- docs/onnx-artifacts/PROVISIONING.md
- docs/onnx-artifacts/README.md
- benchmark-results/fd-onnx-workflow-input-alignment-m034-s02.txt
- .gsd/DECISIONS.md
