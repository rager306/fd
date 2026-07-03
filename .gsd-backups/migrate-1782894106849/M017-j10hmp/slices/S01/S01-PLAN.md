# S01: Tagged ONNX 512 legal quality gate

**Goal:** Run the full legal retrieval evaluator against tagged Go ONNX configured with max sequence length 512 and isolated cache namespace.
**Demo:** After this, there is a measured full legal retrieval artifact for tagged Go ONNX at max sequence length 512, with isolated cache namespace and no raw text leaks.

## Must-Haves

- Required local ONNX/native tokenizer artifacts are present.
- Tagged ONNX service starts with `ONNX_MAX_SEQUENCE_LENGTH=512` and isolated `EMBEDDING_CACHE_VERSION`.
- Legal retrieval evaluator runs against full corpus.
- Artifact is saved under `benchmark-results/` and contains no raw legal text.
- Service cleanup is verified.

## Proof Level

- This slice proves: Live TEI and tagged ONNX services plus evaluator artifact and hygiene checks.

## Integration Closure

Provides measured 512-token Go runtime quality evidence for S02 decision.

## Verification

- Captures runtime command/config and quality metrics in benchmark artifact.

## Tasks

- [x] **T01: Prepare 512 gate command plan** `est:small`
  Inspect existing evaluator CLI and current runtime prerequisites for running tagged Go ONNX at max sequence length 512. Confirm required local artifacts exist and identify the exact evaluator command.
  - Verify: Command plan identifies API URLs, namespace, sequence length, and output artifact.

- [x] **T02: Start tagged ONNX 512 service** `est:small`
  Start tagged Go ONNX service with `ONNX_MAX_SEQUENCE_LENGTH=512`, isolated Redis namespace, and native HF tokenizer; verify health.
  - Verify: `/health` returns ok for the tagged ONNX service.

- [x] **T03: Run full legal quality gate** `est:medium`
  Run the legal retrieval evaluator against TEI and tagged ONNX 512, then check artifact hygiene and summarize metrics.
  - Files: `benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt`
  - Verify: Evaluator exits 0 or records explicit fail verdict; artifact exists and raw legal text leak check passes.

- [x] **T04: Cleanup runtime** `est:small`
  Stop tagged ONNX service and verify no stale benchmark runtime remains.
  - Verify: Background process list shows the tagged ONNX service stopped.

## Files Likely Touched

- benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt
