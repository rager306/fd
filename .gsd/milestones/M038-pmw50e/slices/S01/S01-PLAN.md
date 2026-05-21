# S01: Go ONNX runtime smoke proof

**Goal:** Run local Go ONNX target-runtime smoke proof for current artifact.
**Demo:** After this, the current ONNX artifact has fresh Go live embedder/API smoke evidence with isolated namespace.

## Must-Haves

- Tagged Go live ONNX embedder test passes.
- Go ONNX API starts on local port with isolated cache namespace.
- `/health` exposes safe runtime metadata.
- `/v1/embeddings` returns 1024-dimensional embedding.
- Outcome records no production promotion.

## Proof Level

- This slice proves: Live Go tagged test plus local Go API smoke; no external action.

## Integration Closure

Consumes M037 target-runtime contract and existing M010/M018/M019 artifact/runtime evidence.

## Verification

- Captures health metadata, dimensions, cache namespace, and command evidence.

## Tasks

- [x] **T01: Check Go runtime prerequisites** `est:small`
  Check local prerequisites for Go target-runtime smoke: ONNX artifact, native tokenizer, tokenizer JSON, ONNX Runtime shared library, Redis availability, and clean local ports. Do not start services yet except read-only probes.
  - Verify: Prerequisite check reports exact availability or blocker without leaking secrets.

- [x] **T02: Run live Go ONNX embedder test** `est:small`
  Run live tagged Go ONNX embedder test against the current local artifact and native HF tokenizer path.
  - Verify: `TestONNXEmbedderLiveLocalArtifact` passes with ONNX/HF tokenizer tags and configured runtime library.

- [x] **T03: Run Go ONNX API smoke** `est:medium`
  Start local Go ONNX API with isolated cache namespace on port 18000, verify `/health` runtime metadata and `/v1/embeddings` 1024-dimensional response, then stop the server.
  - Files: `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt`
  - Verify: Health and embeddings smoke pass; server stopped; outcome has no raw input text/secrets/signed URLs.

- [x] **T04: Verify S01 smoke proof** `est:small`
  Complete S01 with smoke proof summary and GitNexus detect scope.
  - Verify: Smoke outcome checks, no background process leaks, port clean, GitNexus detect low/expected.

## Files Likely Touched

- benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt
