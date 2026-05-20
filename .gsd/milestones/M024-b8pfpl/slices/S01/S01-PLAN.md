# S01: Packaged ONNX performance benchmark

**Goal:** Run benchmark.py against the packaged ONNX Docker image with isolated cache namespace and a target-specific restart command.
**Demo:** After this, there is a packaged ONNX Docker performance artifact with sanitized config and comparable metrics.

## Must-Haves

- Packaged ONNX container starts healthy on port 18000.
- Benchmark command uses `uv run --python 3.13`.
- Benchmark env includes packaged runtime label, build tags, manifests, ONNX runtime library marker, cache namespace, and restart command.
- Benchmark writes artifact under `benchmark-results/`.
- Artifact passes basic hygiene and contains expected config markers.

## Proof Level

- This slice proves: Endpoint health, benchmark artifact, config snapshot, cleanup.

## Integration Closure

Produces packaged performance evidence for the opt-in ONNX Docker image.

## Verification

- Benchmark artifact captures sanitized config, image/runtime labels, manifests, Redis metadata, and restart/L2 evidence.

## Tasks

- [x] **T01: Prepare packaged ONNX benchmark target** `est:small`
  Prepare the packaged ONNX benchmark target: ensure TEI/default stack is not modified, start `fd-api:onnx1024-m022-final` on port 18000 with cache namespace `m024-onnx-docker-benchmark`, verify `/health` and embedding dimensions, and confirm restart command strategy.
  - Verify: Health and smoke embedding pass; restart command targets packaged ONNX container.

- [x] **T02: Run packaged ONNX benchmark** `est:medium`
  Run `benchmark.py` against packaged ONNX Docker image and write `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt` with sanitized effective config.
  - Files: `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt`
  - Verify: Benchmark exits 0 and artifact contains packaged runtime config markers.

- [x] **T03: Validate benchmark artifact and cleanup** `est:small`
  Validate benchmark artifact hygiene, extract key metrics, and clean packaged ONNX container/port.
  - Files: `benchmark-results/fd-benchmark-m024-onnx-docker1024.txt`
  - Verify: Artifact exists, expected markers present, no forbidden tracked binaries, no background processes, port 18000 clean.

## Files Likely Touched

- benchmark-results/fd-benchmark-m024-onnx-docker1024.txt
