# S01: Packaged ONNX smoke proof

**Goal:** Build and smoke-test packaged Go ONNX Docker runtime.
**Demo:** After this, a dedicated packaged ONNX Docker image has been built or verified and smoke-tested through its Go API.

## Must-Haves

- ONNX image build or existing image verification succeeds.
- Container starts on port 18000 with isolated namespace.
- `/health` reports verified ONNX runtime metadata.
- `/v1/embeddings` returns 1024-dimensional normalized output.
- Container stopped and port clean.

## Proof Level

- This slice proves: Docker build/container smoke proof.

## Integration Closure

Consumes M038 local Go evidence and existing Dockerfile.onnx packaging contract.

## Verification

- Captures image metadata, health/runtime metadata, dimensions, and namespace.

## Tasks

- [x] **T01: Check packaged runtime prerequisites** `est:small`
  Inspect Docker ONNX build script and Dockerfile contract, then verify local required artifacts/paths before build.
  - Files: `Dockerfile.onnx`, `tools/build_onnx_image.sh`
  - Verify: Build script/artifact prerequisite checks pass.

- [x] **T02: Build packaged ONNX image** `est:medium`
  Build dedicated ONNX Docker image for M039 from current artifacts and record image id/digest-like metadata.
  - Verify: Docker image build succeeds and image id recorded.

- [x] **T03: Run packaged ONNX smoke** `est:medium`
  Run packaged ONNX container on port 18000 with isolated namespace, verify `/health` and `/v1/embeddings`, write smoke artifact, then stop container.
  - Files: `benchmark-results/fd-onnx-docker-smoke-m039-s01.txt`
  - Verify: Packaged health/embedding smoke passes; no raw text/secrets; container stopped; port clean.

- [x] **T04: Verify packaged smoke proof** `est:small`
  Verify S01 scope, background/container cleanup, port cleanliness, and GitNexus detect.
  - Verify: Outcome checks, container/process cleanup, port clean, GitNexus detect pass.

## Files Likely Touched

- Dockerfile.onnx
- tools/build_onnx_image.sh
- benchmark-results/fd-onnx-docker-smoke-m039-s01.txt
