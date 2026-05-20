# S01: Packaged ONNX legal quality gate

**Goal:** Run the Russian/legal retrieval gate against TEI baseline and the packaged ONNX Docker image with isolated cache namespace.
**Demo:** After this, the packaged ONNX Docker image has a Russian/legal retrieval gate artifact or a concrete environment blocker.

## Must-Haves

- TEI baseline is healthy or a concrete blocker is recorded.
- Packaged ONNX Docker endpoint is healthy on port 18000.
- Evaluator runs with isolated ONNX cache namespace.
- Artifact has sanitized config and no raw legal text.
- Strict gate passes or failure/blocker is recorded.

## Proof Level

- This slice proves: Container health, evaluator artifact, pass/fail metrics, raw text leak check.

## Integration Closure

Proves or blocks packaged ONNX legal quality before performance or production rollout work.

## Verification

- Produces sanitized legal-quality artifact and records endpoint/runtime/cache labels.

## Tasks

- [x] **T01: Prepare packaged ONNX legal gate environment** `est:medium`
  Prepare packaged legal gate environment: verify TEI baseline health, build or reuse the M022 ONNX image, start packaged ONNX on port 18000 with an isolated cache namespace, and smoke `/health` plus non-legal embedding dimensions.
  - Verify: TEI `/health` and ONNX `/health` pass; ONNX smoke returns 1024 dimensions; port/process state is known.

- [x] **T02: Run packaged legal retrieval evaluator** `est:medium`
  Run `tools/evaluate_legal_retrieval.py` against TEI at 8000 and packaged ONNX at 18000, writing a M023 artifact under `benchmark-results/` with runtime labels and cache namespaces.
  - Files: `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt`
  - Verify: Evaluator exits 0 for pass or nonzero with blocked/fail artifact; artifact excludes raw legal text.

- [x] **T03: Validate artifact hygiene and cleanup** `est:small`
  Validate legal artifact hygiene and clean packaged runtime: check no raw corpus lines leak, stop ONNX container, confirm port 18000 clean, and record pass/fail outcome.
  - Files: `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt`
  - Verify: Raw text leak check passes; no background processes; port 18000 clean.

## Files Likely Touched

- benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt
