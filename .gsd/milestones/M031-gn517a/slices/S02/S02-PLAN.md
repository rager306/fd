# S02: Source contract documentation and closure

**Goal:** Persist source contract updates and close M031 locally.
**Demo:** After this, manifests/docs/outcome record the source contract and remaining blockers, with final guardrails passing.

## Must-Haves

- Provisioning docs and relevant manifests reflect source status.
- Decision recorded.
- Hosted workflow remains unrun and blocked until explicit approval and real sources.
- Final checks pass and GitNexus post-reindex detect is clean.

## Proof Level

- This slice proves: Docs/manifests/outcome checks, provisioning dry-run/verifier, default guardrails, commit/reindex.

## Integration Closure

Keeps artifact provisioning docs and decisions aligned before hosted proof.

## Verification

- Updates docs/decision surfaces for future workflow run.

## Tasks

- [x] **T01: Persist source contract in docs and manifests** `est:medium`
  Update ONNX artifact manifests and provisioning docs with pinned source candidates, ONNX model blocker, and runtime/source checksum details without changing TEI/default runtime behavior.
  - Files: `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`, `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`, `docs/onnx-artifacts/PROVISIONING.md`
  - Verify: JSON parses; docs contain source statuses and no production promotion language.

- [x] **T02: Record source contract outcome and decision** `est:small`
  Create M031 outcome artifact and decision summarizing selected immutable candidates, ONNX model blocker, and remaining hosted proof blockers.
  - Files: `benchmark-results/fd-onnx-source-contract-m031-s02.txt`, `.gsd/DECISIONS.md`
  - Verify: Outcome has no raw text/secrets/signed URLs and decision is recorded.

- [x] **T03: Verify and close milestone** `est:medium`
  Run final guardrails, validate and complete M031, checkpoint DB, commit locally, run GitNexus reindex/detect, and report state.
  - Verify: Docs/manifests checks, py_compile/provisioning/verifier allow-missing, Go checks as relevant, actionlint, binary hygiene, GitNexus detect, commit, reindex.

## Files Likely Touched

- docs/onnx-artifacts/user-bge-m3-dense-fp32.json
- docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
- docs/onnx-artifacts/PROVISIONING.md
- benchmark-results/fd-onnx-source-contract-m031-s02.txt
- .gsd/DECISIONS.md
