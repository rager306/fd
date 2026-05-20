# S01: Tagged ONNX 1024 legal quality gate

**Goal:** Run the full legal retrieval evaluator against tagged Go ONNX configured with max sequence length 1024 and isolated cache namespace.
**Demo:** After this, there is a measured full legal retrieval artifact for tagged Go ONNX at max sequence length 1024, with isolated cache namespace and sanitized output.

## Must-Haves

- Required ONNX/native tokenizer/corpus artifacts are present.
- Tagged ONNX service starts with `ONNX_MAX_SEQUENCE_LENGTH=1024` and isolated `EMBEDDING_CACHE_VERSION`.
- Legal retrieval evaluator runs against full corpus.
- Artifact is saved under `benchmark-results/` and contains no raw legal text.
- Service cleanup is verified.

## Proof Level

- This slice proves: Live tagged ONNX service, evaluator artifact, hygiene check, cleanup.

## Integration Closure

Provides measured evidence on whether 1024 alone can pass strict legal vector equivalence.

## Verification

- Captures 1024 runtime config and legal quality metrics in benchmark artifact.

## Tasks

- [x] **T01: Prepare 1024 gate command plan** `est:small`
  Confirm runtime prerequisites and exact command for the tagged Go ONNX 1024 legal gate.
  - Verify: Command plan identifies API URLs, namespace, sequence length, and output artifact.

- [x] **T02: Start tagged ONNX 1024 service** `est:small`
  Start tagged Go ONNX service with max sequence length 1024, isolated Redis namespace, and native HF tokenizer; verify health.
  - Verify: `/health` returns ok for the tagged ONNX service.

- [x] **T03: Run full legal quality gate** `est:medium`
  Run the legal retrieval evaluator against TEI and tagged ONNX 1024, then check artifact hygiene and summarize metrics.
  - Files: `benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt`
  - Verify: Evaluator exits 0 or records explicit fail verdict; artifact exists and raw legal text leak check passes.

- [x] **T04: Cleanup runtime** `est:small`
  Stop tagged ONNX service and verify no stale benchmark runtime remains.
  - Verify: Background process list shows the tagged ONNX service stopped.

## Files Likely Touched

- benchmark-results/fd-legal-retrieval-m018-s01-onnx1024.txt
