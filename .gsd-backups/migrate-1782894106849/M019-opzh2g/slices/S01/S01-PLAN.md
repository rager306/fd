# S01: ONNX 1024 performance benchmark

**Goal:** Run benchmark.py against tagged Go ONNX 1024 and compare the result to TEI baseline evidence.
**Demo:** After this, ONNX 1024 has a measured benchmark artifact comparable to the TEI baseline with sanitized config and isolated namespace.

## Must-Haves

- Benchmark command plan is explicit.
- Tagged ONNX 1024 service starts and health passes.
- Benchmark writes `benchmark-results/fd-benchmark-m019-onnx1024.txt`.
- Benchmark artifact records sanitized effective config, including sequence length and cache namespace.
- Runtime cleanup is verified.

## Proof Level

- This slice proves: Live benchmark run, artifact hygiene, runtime cleanup.

## Integration Closure

Produces performance evidence for S02 decision on whether to proceed to packaging or tune first.

## Verification

- Captures ONNX 1024 benchmark artifact with effective config and runtime metadata.

## Tasks

- [x] **T01: Prepare benchmark command plan** `est:small`
  Inspect benchmark.py options and existing M014 artifacts to build the exact ONNX 1024 benchmark command with isolated namespace and restart command.
  - Verify: Command plan identifies API_URL, output artifact, runtime label, namespace, and restart command.

- [x] **T02: Start ONNX 1024 benchmark service** `est:small`
  Start tagged Go ONNX service with `ONNX_MAX_SEQUENCE_LENGTH=1024`, isolated namespace, and native HF tokenizer for benchmarking.
  - Verify: `/health` returns ok for the benchmark service.

- [x] **T03: Run ONNX 1024 benchmark** `est:medium`
  Run `benchmark.py` against ONNX 1024 with uv/Python 3.13 and save the artifact.
  - Files: `benchmark-results/fd-benchmark-m019-onnx1024.txt`
  - Verify: Benchmark exits 0 and artifact hygiene check passes.

- [x] **T04: Cleanup benchmark runtime** `est:small`
  Stop the tagged ONNX benchmark service and verify no stale runtime remains.
  - Verify: Background process list shows no benchmark service.

## Files Likely Touched

- benchmark-results/fd-benchmark-m019-onnx1024.txt
